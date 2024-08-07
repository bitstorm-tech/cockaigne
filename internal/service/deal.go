package service

import (
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/bitstorm-tech/cockaigne/internal/model"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func GetCategories() []model.Category {
	var categories []model.Category
	err := persistence.DB.Select(&categories, "select * from categories where is_active = true")
	if err != nil {
		zap.L().Sugar().Errorf("can't get categories: %v", err)
	}

	return categories
}

func GetCategory(id int) (model.Category, error) {
	var category model.Category
	err := persistence.DB.Get(&category, "select * from categories where id = $1", id)
	if err != nil {
		return model.Category{}, fmt.Errorf("can't get category (id=%d): %v", id, err)
	}

	return category, nil
}

func SaveDeal(deal model.Deal) (dealId uuid.UUID, templateId uuid.UUID, err error) {
	err = persistence.DB.Get(&dealId,
		"insert into deals (dealer_id, title, description, category_id, duration_in_hours, start, start_instantly, own_end_date, payment_state, template) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, false) returning id",
		deal.DealerId,
		deal.Title,
		deal.Description,
		deal.CategoryId,
		deal.DurationInHours,
		deal.Start,
		deal.StartInstantly,
		deal.OwnEndDate,
		model.DealPaymentStatePending,
	)

	if err != nil {
		return dealId, templateId, err
	}

	if deal.IsTemplate {
		err = persistence.DB.Get(
			&templateId,
			"insert into deals (dealer_id, title, description, category_id, duration_in_hours, start, start_instantly, own_end_date, template) values ($1, $2, $3, $4, $5, $6, $7, $8, true) returning id",
			deal.DealerId,
			deal.Title,
			deal.Description,
			deal.CategoryId,
			deal.DurationInHours,
			deal.Start,
			deal.StartInstantly,
			deal.OwnEndDate,
		)
	}

	return dealId, templateId, err
}

func GetDeal(id string) (model.Deal, error) {
	var deal model.Deal
	err := persistence.DB.Get(&deal, "select * from deals where id = $1", id)
	if err != nil {
		return model.Deal{}, fmt.Errorf("can't get deal from database: %v", err)
	}

	return deal, nil
}

type SpatialDealFilter interface {
	ToGeometry() (string, error)
}

type BoundingBoxDealFilter struct {
	BoundingBox string
}

type RadiusDealFilter struct {
	Point  model.Point
	Radius int
}

func (filter BoundingBoxDealFilter) ToGeometry() (string, error) {
	if len(filter.BoundingBox) == 0 {
		return "", fmt.Errorf("BoundingBoxDealFilter needs a valid bounding box")
	}

	coordinates := strings.Split(filter.BoundingBox, ",")
	return fmt.Sprintf(
		"ST_MakeEnvelope(%s, %s, %s, %s, 4326)",
		coordinates[0],
		coordinates[1],
		coordinates[2],
		coordinates[3],
	), nil
}

func (filter RadiusDealFilter) ToGeometry() (string, error) {
	if filter.Point.Lat == 0.0 || filter.Point.Lon == 0.0 || filter.Radius == 0 {
		return "", fmt.Errorf("RadiusDealFilter needs a valid point and raidus")
	}

	pointString := fmt.Sprintf("%f,%f", filter.Point.Lon, filter.Point.Lat)

	return fmt.Sprintf("ST_Buffer(ST_MakePoint(%s)::geography, %d)::geometry", pointString, filter.Radius), nil
}

func CreateSpatialDealFilter(userId string) (SpatialDealFilter, error) {
	account, err := GetAccount(userId)
	if err != nil {
		return nil, err
	}

	if !account.Location.Valid {
		return nil, fmt.Errorf("can't create SpatialDealFilter -> account has no location")
	}

	if account.SearchRadiusInMeters == 0 {
		return nil, fmt.Errorf("can't create SpatialDealFilter -> account has no search radius")
	}

	point, err := model.NewPointFromString(account.Location.String)
	if err != nil {
		return nil, err
	}

	radius := account.SearchRadiusInMeters

	filter := RadiusDealFilter{
		Point:  point,
		Radius: radius,
	}

	return filter, nil
}

func GetActiveDeals(filter SpatialDealFilter, user User) ([]model.DealView, error) {
	query := "select *, public.st_x(location) || ',' || public.st_y(location) as location from active_deals_view"

	geom, err := filter.ToGeometry()
	if err != nil {
		return []model.DealView{}, err
	}

	query += fmt.Sprintf(" where ST_Within(location, %s)", geom)

	addCategoryIdFilter(user, &query)

	var deals []model.DealView
	err = persistence.DB.Select(&deals, query)

	if err != nil {
		return []model.DealView{}, fmt.Errorf("can't get active deals: %v", err)
	}

	return deals, nil
}

func addCategoryIdFilter(user User, query *string) {
	var favoriteCategoryIds []int

	if user.IsBasicUser {
		basicUserFilter := GetBasicUserFilter(user.ID.String())
		favoriteCategoryIds = basicUserFilter.SelectedCategories
	} else {
		favoriteCategoryIds = GetFavoriteCategoryIds(user.ID)
	}

	if len(favoriteCategoryIds) == 0 {
		return
	}

	if strings.Contains(strings.ToLower(*query), "where") {
		*query += " and "
	} else {
		*query += " where "
	}

	categoryIdsString := fmt.Sprintf("%v", favoriteCategoryIds)
	categoryIdsString = strings.Replace(categoryIdsString, " ", ",", -1)

	*query += fmt.Sprintf("category_id = any(array%s)", categoryIdsString)
}

func addSpatialFilterToQuery(user User, query *string) error {
	var filter SpatialDealFilter
	var err error
	if user.IsBasicUser {
		filter = GetBasicUserSpatialFilter(user.ID.String())
	} else {
		filter, err = CreateSpatialDealFilter(user.ID.String())
		if err != nil {
			zap.L().Sugar().Error("can't create SpatialDealFilter: ", err)
		}
	}

	geom, err := filter.ToGeometry()
	if err != nil {
		return err
	}

	if strings.Contains(strings.ToLower(*query), "where") {
		*query += " and "
	} else {
		*query += " where "
	}

	*query += fmt.Sprintf("ST_Within(location, %s)", geom)
	return nil
}

func GetDealHeaders(state model.DealState, user *User, dealerId string) (model.DealHeaders, error) {
	if state != model.DealStateFuture && state != model.DealStateActive && state != model.DealStatePast && state != model.DealStateTemplate {
		return model.DealHeaders{}, fmt.Errorf("unknown deal state: %s", state)
	}

	query := "select id, title, username, dealer_id, category_id, start_time, start + interval '60 minute' >= now() as CanEdit from active_deals_view"
	switch state {
	case model.DealStatePast:
		query = "select id, title, username, dealer_id, category_id, start_time, start + interval '60 minute' >= now() as CanEdit from past_deals_view"
	case model.DealStateFuture:
		query = "select id, title, username, dealer_id, category_id, start_time, start + interval '60 minute' >= now() as CanEdit from future_deals_view"
	case model.DealStateTemplate:
		query = "select d.id, d.title, a.username, d.dealer_id, d.category_id, d.start + interval '60 minute' >= now() as CanEdit from deals d join accounts a on a.id = d.dealer_id where template = true"
	}

	if len(dealerId) > 0 {
		query += fmt.Sprintf(" where dealer_id = '%s'", dealerId)
	}

	if len(dealerId) == 0 && user != nil {
		err := addSpatialFilterToQuery(*user, &query)
		if err != nil {
			zap.L().Sugar().Error("can't add spatial filter to query: ", err)
		}

		addCategoryIdFilter(*user, &query)
	}

	headers := model.DealHeaders{}
	err := persistence.DB.Select(&headers, query)

	if err != nil {
		return model.DealHeaders{}, fmt.Errorf("can't get active deals: %v", err)
	}

	return headers.RotateByTime(), nil
}

func GetFavoriteDealHeaders(userId string) (model.DealHeaders, error) {
	var headers model.DealHeaders
	err := persistence.DB.Select(
		&headers,
		"select id, dealer_id, title, username, category_id from active_deals_view d join favorite_deals f on d.id = f.deal_id where f.user_id = $1",
		userId,
	)
	if err != nil {
		return model.DealHeaders{}, err
	}

	return headers.RotateByTime(), nil
}

type dealDetailsResult struct {
	model.DealDetails
	DurationInHours int `db:"duration_in_hours"`
	Start           time.Time
}

func GetDealDetails(dealId string) (model.DealDetails, error) {
	var result dealDetailsResult
	err := persistence.DB.Get(&result, "select id, title, description, start, duration_in_hours from deals where id = $1", dealId)
	if err != nil {
		return model.DealDetails{}, fmt.Errorf("can't get deal details of deal %s: %v", dealId, err)
	}

	startAndEndDate := CalculateStartAndEndAsHumanReadable(result.Start, result.DurationInHours)

	return model.DealDetails{
		ID:          result.ID,
		Title:       result.Title,
		Description: result.Description,
		Start:       startAndEndDate.Start,
		End:         startAndEndDate.End,
	}, nil
}

func GetDealReport(dealId string, reporterId string) (model.DealReport, error) {
	reason := ""
	err := persistence.DB.Get(
		&reason,
		"select reason from reported_deals where deal_id = $1 and reporter_id = $2",
		dealId,
		reporterId,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return model.DealReport{}, fmt.Errorf("can't get reason for deal report of deal %s: %v", dealId, err)
	}

	title := ""
	err = persistence.DB.Get(&title, "select title from deals where id = $1", dealId)
	if err != nil {
		return model.DealReport{}, fmt.Errorf("can't get title for deal report of deal %s: %v", dealId, err)
	}

	return model.DealReport{
		Title:  title,
		Reason: reason,
	}, nil
}

func GetDealLikes(dealId string) int {
	likes := 0
	err := persistence.DB.Get(
		&likes,
		"select coalesce(likecount, 0) as likes from like_counts_view where deal_id = $1",
		dealId,
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		zap.L().Sugar().Errorf("can't get like count: %v", err)
		return 0
	}

	return likes
}

func ToggleLikes(dealId string, userId string) int {
	count := 0
	err := persistence.DB.Get(&count, "select count(*)  from likes where deal_id = $1 and user_id = $2", dealId, userId)

	if err != nil {
		zap.L().Sugar().Errorf("can't check if like is already persisted: %v", err)
		return 0
	}

	query := "delete from likes where deal_id = $1 and user_id = $2"
	if count == 0 {
		query = "insert into likes (deal_id, user_id) values ($1, $2)"
	}

	_, err = persistence.DB.Exec(query, dealId, userId)
	if err != nil {
		zap.L().Sugar().Errorf("can't toggle like: %v", err)
		return 0
	}

	likes := 0
	err = persistence.DB.Get(&likes, "select likecount from like_counts_view where deal_id = $1", dealId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		zap.L().Sugar().Errorf("can't get like count for deal %s: %v", dealId, err)
		return 0
	}

	return likes
}

func IsDealLiked(dealId string, userId string) bool {
	isLiked := false
	err := persistence.DB.Get(
		&isLiked,
		"select exists(select user_id from likes where deal_id = $1 and user_id = $2)",
		dealId,
		userId,
	)
	if err != nil {
		zap.L().Sugar().Errorf("can't check if user has liked the deal %s: %v", dealId, err)
		return false
	}

	return isLiked
}

func IsDealFavorite(dealId string, userId string) bool {
	isFavorite := false
	err := persistence.DB.Get(
		&isFavorite,
		"select exists(select user_id from favorite_deals where deal_id = $1 and user_id = $2)",
		dealId,
		userId,
	)
	if err != nil {
		zap.L().Sugar().Errorf("can't check if deal %s is favorite: %v", dealId, err)
		return false
	}

	return isFavorite
}

func ToggleFavorite(dealId string, userId string) bool {
	isFavorite := IsDealFavorite(dealId, userId)

	var err error
	if isFavorite {
		_, err = persistence.DB.Exec("delete from favorite_deals where deal_id = $1 and user_id = $2", dealId, userId)
	} else {
		_, err = persistence.DB.Exec("insert into favorite_deals (user_id, deal_id) values ($1, $2)", userId, dealId)
	}

	if err != nil {
		zap.L().Sugar().Errorf("can't check if deal %s is favorite: %v", dealId, err)
		return false
	}

	return !isFavorite
}

func UploadDealImage(image *multipart.FileHeader, dealId string, name string) error {
	_, err := persistence.UploadDealImage(dealId, name, image)

	return err
}

func GetDealImageUrls(dealId string) ([]string, error) {
	return persistence.GetDealImageUrls(dealId)
}

func SaveDealReport(dealId string, reporterId string, reason string) error {
	_, err := persistence.DB.Exec(
		"insert into reported_deals (reporter_id, deal_id, reason) values ($1, $2, $3)",
		reporterId,
		dealId,
		reason,
	)

	return err
}

func RemoveDealFavorite(dealId string, userId string) error {
	_, err := persistence.DB.Exec(
		"delete from favorite_deals where deal_id = $1 and user_id = $2",
		dealId,
		userId,
	)

	return err
}

func GetFavoriteDealerDealHeaders(userId string) (model.DealHeaders, error) {
	var header model.DealHeaders

	err := persistence.DB.Select(
		&header,
		"select id, d.dealer_id, title, username, category_id from active_deals_view d join favorite_dealers f on d.dealer_id = f.dealer_id where user_id = $1 order by start_time",
		userId,
	)
	if err != nil {
		return nil, err
	}

	return header.RotateByTime(), nil
}

func GetTopDealHeaders(userId string, limit int) ([]model.DealHeader, error) {
	filter, err := CreateSpatialDealFilter(userId)
	if err != nil {
		return []model.DealHeader{}, err
	}

	filterGeom, err := filter.ToGeometry()
	if err != nil {
		return []model.DealHeader{}, err
	}

	query := fmt.Sprintf(
		"select id, dealer_id, title, username, category_id from top_deals_view where st_within(location, %s) limit $1",
		filterGeom,
	)

	var header []model.DealHeader
	err = persistence.DB.Select(
		&header,
		query,
		limit,
	)

	return header, err
}

func GetFavoriteDealsCount(userId string) (int, error) {
	count := 0
	err := persistence.DB.Get(
		&count,
		"select count(*) from favorite_deals f join active_deals_view a on a.id = f.deal_id where f.user_id = $1",
		userId,
	)

	return count, err
}

func HasActiveSubscription(dealerId string) (bool, error) {
	hasActiveSub := false
	err := persistence.DB.Get(
		&hasActiveSub,
		"select exists(select * from subscriptions where state = $1 and account_id = $2)",
		model.SubscriptionStateActive,
		dealerId,
	)

	return hasActiveSub, err
}

func GetHighestVoucherDiscount(dealerId string) (int, error) {
	highestDiscount := 0
	err := persistence.DB.Get(
		&highestDiscount,
		"select coalesce(max(discount_in_percent), 0) from active_vouchers_view where account_id = $1",
		dealerId,
	)

	return highestDiscount, err
}

func FormatPrice(price float64) string {
	priceString := fmt.Sprintf("%f", price)
	dotIndex := strings.Index(priceString, ".")

	return priceString[:dotIndex+3]
}

func FormatPriceWithDiscount(price float64, discountInPercent int) string {
	percent := (100.0 - float64(discountInPercent)) / 100.0
	priceWithDiscount := price * percent

	return FormatPrice(priceWithDiscount)
}

type StartAndEndDate struct {
	Start string
	End   string
}

func CalculateStartAndEndAsHumanReadable(start time.Time, durationInHours int) StartAndEndDate {
	return StartAndEndDate{
		Start: start.Format("02.01.2006 15:04"),
		End:   start.Add(time.Duration(durationInHours) * time.Hour).Format("02.01.2006 15:04"),
	}
}

func NewDealsAvailable(user User, oldDealIds []string) (bool, error) {
	var filter SpatialDealFilter
	var favoriteCategoryIds []int
	var err error

	if user.IsBasicUser {
		filter = GetBasicUserSpatialFilter(user.ID.String())
		favoriteCategoryIds = GetBasicUserFilter(user.ID.String()).SelectedCategories
	} else {
		favoriteCategoryIds = GetFavoriteCategoryIds(user.ID)
		filter, err = CreateSpatialDealFilter(user.ID.String())
		if err != nil {
			return false, err
		}
	}

	searchRadiusFilterGeometry, err := filter.ToGeometry()
	if err != nil {
		return false, err
	}

	query := fmt.Sprintf("select count(*) from active_deals_view where st_within(location, %s)", searchRadiusFilterGeometry)

	var params []interface{}
	if len(oldDealIds) > 0 {
		query += " and id not in (?)"
		query, params, err = sqlx.In(query, oldDealIds)
		if err != nil {
			return false, err
		}
	}

	if len(favoriteCategoryIds) > 0 {
		categoryIdsString := fmt.Sprintf("%v", favoriteCategoryIds)
		categoryIdsString = strings.Replace(categoryIdsString, " ", ",", -1)
		query += fmt.Sprintf(" and category_id = any(array%s)", categoryIdsString)
	}

	query = persistence.DB.Rebind(query)

	var newDealsAvailable int
	err = persistence.DB.Get(&newDealsAvailable, query, params...)
	if err != nil {
		return false, err
	}

	return newDealsAvailable > 0, nil
}

func MarkDealAsPayed(dealId string) error {
	_, err := persistence.DB.Exec(
		"update deals set payment_state = $1 where id = $2",
		model.DealPaymentStatePayed,
		dealId,
	)

	return err
}

func DealClicked(userId, dealId string) error {
	_, err := persistence.DB.Exec(
		"insert into deal_clicks (deal_id, account_id) values ($1, $2) on conflict do nothing",
		dealId,
		userId,
	)

	return err
}

func GetDealStatistics(dealId string) (model.DealStatistics, error) {
	var statistics model.DealStatistics
	err := persistence.DB.Get(
		&statistics,
		"select * from statistics_view where deal_id=$1",
		dealId,
	)

	return statistics, err
}

func UpdateDeal(deal model.Deal) error {
	var err error

	if deal.IsTemplate {
		_, err = persistence.DB.Exec(
			"update deals set start_instantly = $1, own_end_date = $2, start = $3, category_id = $4, title = $5, description = $6, duration_in_hours = $7 where id = $8 and dealer_id = $9",
			deal.StartInstantly,
			deal.OwnEndDate,
			deal.Start,
			deal.CategoryId,
			deal.Title,
			deal.Description,
			deal.DurationInHours,
			deal.ID,
			deal.DealerId,
		)
	} else {
		_, err = persistence.DB.Exec(
			"update deals set category_id = $1, title = $2, description = $3 where id = $4 and dealer_id = $5",
			deal.CategoryId,
			deal.Title,
			deal.Description,
			deal.ID,
			deal.DealerId,
		)
	}

	return err
}

func DeleteDealImage(dealId string, imageName string) error {
	return persistence.DeleteDealImage(dealId, imageName)
}

func CopyDealImages(fromDealId string, toDealId string) error {
	return persistence.CopyDealImages(fromDealId, toDealId)
}
