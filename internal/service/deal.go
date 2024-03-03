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

func SaveDeal(deal model.Deal) (uuid.UUID, error) {
	var dealId uuid.UUID
	err := persistence.DB.Get(&dealId,
		"insert into deals (dealer_id, title, description, category_id, duration_in_hours, start, template) values ($1, $2, $3, $4, $5, $6, false) returning id",
		deal.DealerId,
		deal.Title,
		deal.Description,
		deal.CategoryId,
		deal.DurationInHours,
		deal.Start,
	)

	if err != nil {
		return uuid.UUID{}, err
	}

	if deal.IsTemplate {
		_, err = persistence.DB.Exec(
			"insert into deals (dealer_id, title, description, category_id, duration_in_hours, start, template) values ($1, $2, $3, $4, $5, $6, true) returning id",
			deal.DealerId,
			deal.Title,
			deal.Description,
			deal.CategoryId,
			deal.DurationInHours,
			deal.Start,
		)
	}

	if err != nil {
		return uuid.UUID{}, err
	}

	return dealId, nil
}

func GetDeal(id string) (model.Deal, error) {
	var deal model.Deal
	err := persistence.DB.Get(&deal, "select * from deals where id = $1", id)
	if err != nil {
		return model.Deal{}, fmt.Errorf("can't get deal from database: %v", err)
	}

	return deal, nil
}

type SpartialDealFilter interface {
	ToGeometry() (string, error)
}

type BoundingBoxDealFilter struct {
	BoundingBox string
}

func (filter BoundingBoxDealFilter) ToGeometry() (string, error) {
	if len(filter.BoundingBox) == 0 {
		return "", fmt.Errorf("BoundingBoxDealFilter needs a valid bounding box")
	}

	coordinates := strings.Split(filter.BoundingBox, ",")
	return fmt.Sprintf(
		"ST_Envelope('LINESTRING(%s %s, %s %s)'::geography::geometry)",
		coordinates[0],
		coordinates[1],
		coordinates[2],
		coordinates[3],
	), nil
}

func GetDealsFromView(state model.State, filter SpartialDealFilter, dealerId *string) ([]model.DealView, error) {
	if state != model.Future && state != model.Active && state != model.Past {
		return []model.DealView{}, fmt.Errorf("unknown deal state: %s", state)
	}

	statement := "select *, public.st_x(location) || ',' || public.st_y(location) as location from active_deals_view"
	switch state {
	case model.Past:
		statement = "select *, public.st_x(location) || ',' || public.st_y(location) as location from past_deals_view"
	case model.Future:
		statement = "select *, public.st_x(location) || ',' || public.st_y(location) as location from future_deals_view"
	}

	if dealerId != nil {
		statement += fmt.Sprintf(" where dealer_id = '%s'", *dealerId)
	}

	if filter != nil {
		if dealerId != nil {
			statement += " and "
		} else {
			statement += " where "
		}

		geom, err := filter.ToGeometry()
		if err != nil {
			return []model.DealView{}, fmt.Errorf("can't convert filter to valid geometry: %+v", err)
		}

		statement += fmt.Sprintf("st_within(location, %s)", geom)
	}

	var deals []model.DealView
	err := persistence.DB.Select(&deals, statement)

	if err != nil {
		return []model.DealView{}, fmt.Errorf("can't get active deals: %v", err)
	}

	return deals, nil
}

func GetDealHeaders(state model.State, dealerId string) (model.DealHeaders, error) {
	if state != model.Future && state != model.Active && state != model.Past && state != model.Template {
		return model.DealHeaders{}, fmt.Errorf("unknown deal state: %s", state)
	}

	statement := "select id, title, username, dealer_id, category_id, start_time from active_deals_view"
	switch state {
	case model.Past:
		statement = "select id, title, username, dealer_id, category_id, start_time from past_deals_view"
	case model.Future:
		statement = "select id, title, username, dealer_id, category_id, start_time from future_deals_view"
	case model.Template:
		statement = "select d.id, d.title, a.username, d.dealer_id, d.category_id from deals d join accounts a on a.id = d.dealer_id where template = true"
	}

	if len(dealerId) > 0 {
		statement += fmt.Sprintf(" where dealer_id = '%s'", dealerId)
	}

	var headers model.DealHeaders
	err := persistence.DB.Select(&headers, statement)

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

func GetTemplateHeaders(dealerId string) ([]model.Deal, error) {
	var templates []model.Deal
	err := persistence.DB.Select(&templates, "select * from deals where template = true and dealer_id = $1", dealerId)
	if err != nil {
		return []model.Deal{}, fmt.Errorf("can't get templates: %v", err)
	}

	return templates, nil
}

func UploadDealImage(image *multipart.FileHeader, dealId string, prefix string) error {
	tokens := strings.Split(image.Filename, ".")
	fileExtension := tokens[len(tokens)-1]
	path := fmt.Sprintf("%s/%s/%s%d.%s", persistence.DealImagesFolder, dealId, prefix, time.Now().UnixMilli(), fileExtension)

	return persistence.UploadImage(path, image)
}

func GetDealImageUrls(dealId string) ([]string, error) {
	path := fmt.Sprintf("%s/%s", persistence.DealImagesFolder, dealId)

	return persistence.GetImageUrls(path)
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

func GetTopDealHeaders(limit int) ([]model.DealHeader, error) {
	var header []model.DealHeader
	err := persistence.DB.Select(
		&header,
		"select id, dealer_id, title, username, category_id from top_deals_view limit $1",
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
		model.SubActive,
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
		Start: start.Format("02.01.2006 um 15:04"),
		End:   start.Add(time.Duration(durationInHours) * time.Hour).Format("02.01.2006 um 15:04"),
	}
}
