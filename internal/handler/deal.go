package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bitstorm-tech/cockaigne/internal/model"
	"github.com/bitstorm-tech/cockaigne/internal/redirect"
	"github.com/bitstorm-tech/cockaigne/internal/service"
	"github.com/bitstorm-tech/cockaigne/internal/view"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RegisterDealHandlers(e *echo.Echo) {
	e.GET("/deal/:dealId", openDealCreatePage)
	e.GET("/ui/category-select", getCategorySelect)
	e.GET("/deals/:state", getDealList)
	e.GET("/deals-top/:limit", getTopDealsList)
	e.GET("/api/deals", getDealsAsJson)
	e.GET("/deal-statistics/:id", openDealStatisticsPage)
	e.GET("/deal-details/:id", getDealDetails)
	e.GET("/deal-likes/:id", toggleDealLike)
	e.GET("/deal-viewed/:id", dealViewed)
	e.GET("/deal-favorite-button/:id", getDealFavoriteButton)
	e.GET("/deal-favorite-toggle/:id", toggleFavorite)
	e.GET("/deal-favorites-list", getFavoriteDeals)
	e.GET("/deal-image-zoom-modal/:dealId", getImageZoomModal)
	e.GET("/dealer-favorites-list", getFavoriteDealerDeals)
	e.GET("/top-deals", openTopDealsPage)
	e.GET("/ui/deals/report-modal/:id", getReportModal)
	e.GET("/deal-payed/:id", markDealAsPayed)
	e.POST("/deal-new-summary", openNewDealSummaryModal)
	e.POST("/deal-report/:id", saveReport)
	e.POST("/deals", saveDeal)
	e.DELETE("/deal-favorite-remove/:id", removeFavorite)
}

func dealViewed(c echo.Context) error {
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	dealId := c.Param("id")

	if !user.IsBasicUser && !user.IsDealer {
		err := service.DealClicked(user.ID.String(), dealId)
		if err != nil {
			zap.L().Sugar().Error("can't save deal clicked: ", err)
		}
	}

	return nil
}

func openDealStatisticsPage(c echo.Context) error {
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	dealId := c.Param("id")
	dealStatistics, err := service.GetDealStatistics(dealId)
	if err != nil {
		zap.L().Sugar().Error("can't get deal statistics: ", err)
		return view.RenderAlertTranslated("alert.can_t_load_statistics", c)
	}

	if dealStatistics.DealerId.String() != user.ID.String() {
		zap.L().Sugar().Errorf("can't show deal statistics, deal (id=%s) does not belong to dealer (id=%s)", dealStatistics.DealId, user.ID)
		return view.RenderAlertTranslated("alert.not_your_statistics", c)
	}

	return view.Render(view.DealStatistics(dealStatistics), c)
}

func markDealAsPayed(c echo.Context) error {
	dealId := c.Param("id")
	err := service.MarkDealAsPayed(dealId)
	if err != nil {
		zap.L().Sugar().Errorf("can't mark deal (id=%s) as payed: %v", dealId, err)
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func openNewDealSummaryModal(c echo.Context) error {
	dealerId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	title := c.FormValue("title")
	if len(title) == 0 {
		return view.RenderAlertTranslated("alert.enter_title", c)
	}

	description := c.FormValue("description")
	if len(description) == 0 {
		return view.RenderAlertTranslated("alert.enter_description", c)
	}

	timesAndDates, err := calculateDealTimesAndDates(c)
	if err != nil {
		zap.L().Sugar().Error("can't calculate deal times and dates: ", err)
		return view.RenderAlertTranslated("alert.can_t_create_deal", c)
	}

	lang := service.GetLanguageFromCookie(c)

	params := createSubscriptionModalParams(dealerId.String(), timesAndDates)
	if params != nil {
		addStartAndEndToSummaryModalParameter(params, timesAndDates)
		return view.Render(view.NewDealSummaryModal(*params, lang), c)
	}

	params = createDiscountModalParams(dealerId.String(), timesAndDates)
	if params != nil {
		addStartAndEndToSummaryModalParameter(params, timesAndDates)
		return view.Render(view.NewDealSummaryModal(*params, lang), c)
	}

	params = &view.NewDealSummaryModalParameter{
		Err:   false,
		Price: service.FormatPrice(float64(timesAndDates.DurationInDays) * 4.99),
	}

	addStartAndEndToSummaryModalParameter(params, timesAndDates)

	return view.Render(view.NewDealSummaryModal(*params, lang), c)
}

type DealTimesAndDates struct {
	DurationInDays int
	Start          time.Time
	End            time.Time
}

func calculateDealTimesAndDates(c echo.Context) (DealTimesAndDates, error) {
	dealTimesAndDates := DealTimesAndDates{}

	if c.FormValue("startInstantly") == "on" {
		dealTimesAndDates.Start = time.Now()
	} else {
		startDateTime, err := time.Parse("2006-01-02T15:04", c.FormValue("startDate"))
		if err != nil {
			return DealTimesAndDates{}, fmt.Errorf("can't parse start date: %w", err)
		}
		dealTimesAndDates.Start = startDateTime
	}

	if c.FormValue("ownEndDate") == "on" {
		startTime := strings.Split(c.FormValue("startDate"), "T")[1]
		endDate, err := time.Parse("2006-01-0215:04", c.FormValue("endDate")+startTime)
		if err != nil {
			return DealTimesAndDates{}, fmt.Errorf("can't parse end date: %w", err)
		}
		dealTimesAndDates.End = endDate
		dealTimesAndDates.DurationInDays = int(dealTimesAndDates.End.Sub(dealTimesAndDates.Start).Hours()) / 24
	} else {
		durationString := c.FormValue("duration")
		duration, err := strconv.Atoi(durationString)
		if err != nil {
			return DealTimesAndDates{}, fmt.Errorf("can't convert duration from string (%s) to into number: %w", durationString, err)
		}
		dealTimesAndDates.End = dealTimesAndDates.Start.Add(time.Duration(duration*24) * time.Hour)
		dealTimesAndDates.DurationInDays = duration
	}

	return dealTimesAndDates, nil
}

func addStartAndEndToSummaryModalParameter(params *view.NewDealSummaryModalParameter, timesAndDates DealTimesAndDates) {
	params.Start = timesAndDates.Start.Format("02.01.2006 15:04")
	params.End = timesAndDates.End.Format("02.01.2006 15:04")
	params.Duration = fmt.Sprintf("%d", timesAndDates.DurationInDays)
}

func createSubscriptionModalParams(dealerId string, timesAndDates DealTimesAndDates) *view.NewDealSummaryModalParameter {
	params := view.NewDealSummaryModalParameter{}
	hasActiveSub, err := service.HasActiveSubscription(dealerId)
	if err != nil {
		zap.L().Sugar().Error("can't check if dealer has active subscription: ", err)
		params.Err = true
		return &params
	}

	if !hasActiveSub {
		return nil
	}

	freeDaysLeft, err := service.GetFreeDaysLeftFromSubscription(dealerId)
	if err != nil {
		zap.L().Sugar().Error("can't get free days left from subscription: ", err)
		params.Err = true
		return &params
	}

	daysLeftAfterDeal := freeDaysLeft - timesAndDates.DurationInDays
	if daysLeftAfterDeal < 0 {
		params.Price = fmt.Sprintf("%.2f", float64(-1*daysLeftAfterDeal)*4.99)
		params.FreeDaysLeft = "0"
	}

	params.FreeDaysLeft = fmt.Sprintf("%d", daysLeftAfterDeal)
	params.Err = false

	return &params
}

func createDiscountModalParams(dealerId string, timesAndDates DealTimesAndDates) *view.NewDealSummaryModalParameter {
	params := view.NewDealSummaryModalParameter{}
	highestDiscountInPercent, err := service.GetHighestVoucherDiscount(dealerId)
	if err != nil {
		zap.L().Sugar().Error("can't check if dealer has active voucher; ", err)
		params.Err = true
		return &params
	}

	if highestDiscountInPercent == 0 {
		return nil
	}

	price := float64(timesAndDates.DurationInDays) * 4.99
	params.Discount = fmt.Sprintf("%d", highestDiscountInPercent)
	params.Price = service.FormatPrice(price)
	params.PriceWithDiscount = service.FormatPriceWithDiscount(price, highestDiscountInPercent)

	return &params
}

func getTopDealsList(c echo.Context) error {
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	limitString := c.Param("limit")
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		zap.L().Sugar().Error("can't convert limit parameter to int: ", err)
		return view.RenderAlertTranslated("alert.can_t_load_top_deals", c)
	}

	if limit > 100 {
		limit = 100
	}

	header, err := service.GetTopDealHeaders(user.ID.String(), limit)
	if err != nil {
		zap.L().Sugar().Error("can't get top deals: ", err)
		return view.RenderAlertTranslated("alert.can_t_load_top_deals", c)
	}

	return view.Render(view.DealsList(view.DealListTopDeals, user, header, false, view.ActionButtonFavoriteToggle), c)
}

func openTopDealsPage(c echo.Context) error {
	lang := service.GetLanguageFromCookie(c)
	return view.Render(view.TopDealsPage(lang), c)
}

func getFavoriteDealerDeals(c echo.Context) error {
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	headers, err := service.GetFavoriteDealerDealHeaders(user.ID.String())
	if err != nil {
		zap.L().Sugar().Error("can't get favorite dealer deals: ", err)
		return view.RenderAlertTranslated("alert.can_t_load_favorite_deals", c)
	}

	return view.Render(view.DealsList(view.DealListUserFavoriteDealerDeals, user, headers, false, view.ActionButtonNone), c)
}

func getFavoriteDeals(c echo.Context) error {
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	headers, err := service.GetFavoriteDealHeaders(user.ID.String())
	if err != nil {
		zap.L().Sugar().Error("can't get favorite deal headers: ", err)
		return view.RenderAlertTranslated("alert.can_t_load_favorite_deals", c)
	}

	return view.Render(view.DealsList(view.DealListUserFavoriteDeals, user, headers, false, view.ActionButtonRemoveFavorite), c)
}

func openDealCreatePage(c echo.Context) error {
	userId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	dealId := c.Param("dealId")

	var deal model.Deal
	if strings.EqualFold(dealId, "new") {
		deal = model.NewDeal()
		deal.CategoryId = service.GetDefaultCategoryId(userId)
	} else {
		deal, err = service.GetDeal(dealId)
		if err != nil {
			zap.L().Sugar().Errorf("can't get deal %s: %v", dealId, err)
			return view.RenderErrorPageTranslated("alert.can_t_load_deal", c)
		}
	}

	lang := service.GetLanguageFromCookie(c)

	return view.Render(view.DealEdit(deal, lang), c)
}

func getCategorySelect(c echo.Context) error {
	categories := service.GetCategories()
	translationKey := c.QueryParam("translation_key")
	selectedParam := c.QueryParam("selected")
	lang := service.GetLanguageFromCookie(c)
	emptyOptionText := service.T(translationKey, lang)

	selected := -1
	var err error
	if len(selectedParam) > 0 {
		selected, err = strconv.Atoi(c.QueryParam("selected"))
		if err != nil {
			zap.L().Sugar().Error("can't parse selected category: ", err)
			selected = -1
		}
	}

	return view.Render(view.CategorySelect(emptyOptionText, categories, selected), c)
}

func saveDeal(c echo.Context) error {
	userId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	deal, errorMessageTranslationKey := model.DealFromRequest(c)
	if len(errorMessageTranslationKey) > 0 {
		return view.RenderAlertTranslated(errorMessageTranslationKey, c)
	}

	deal.DealerId = userId

	dealId, err := service.SaveDeal(deal)
	if err != nil {
		zap.L().Sugar().Error("can't save deal: ", err)
		return view.RenderAlertTranslated("alert.can_t_save_deal", c)
	}

	form, err := c.MultipartForm()
	if err != nil {
		zap.L().Sugar().Error("can't get multipart form: ", err)
		return view.RenderAlertTranslated("alert.can_t_save_deal", c)
	}

	for index, file := range form.File["images"] {
		err = service.UploadDealImage(file, dealId.String(), fmt.Sprintf("%d", index))
		if err != nil {
			zap.L().Sugar().Error("can't upload deal image: ", err)
			return view.RenderAlertTranslated("alert.can_t_save_deal", c)
		}
	}

	baseUrl := service.GetBaseUrl(c)
	checkoutSession, err := service.DoStripePayment(userId.String(), dealId.String(), deal.DurationInHours/24, baseUrl)
	if err != nil {
		zap.L().Sugar().Error("can't do stripe payment: ", err)
	}

	if checkoutSession != nil {
		c.Response().Header().Set("HX-Redirect", checkoutSession.URL)
	} else {
		err = service.MarkDealAsPayed(dealId.String())
		if err != nil {
			zap.L().Sugar().Errorf("can't mark deal (id=%s) as payed: %v", dealId, err)
		}
		c.Response().Header().Set("HX-Redirect", "/")
	}

	return nil
}

func getDealList(c echo.Context) error {
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	state := model.ToDealState(c.Param("state"))
	dealerId := c.QueryParam("dealer_id")

	headers, err := service.GetDealHeaders(state, &user, dealerId)
	if err != nil {
		zap.L().Sugar().Error("can't get deal headers: ", err)
		return view.RenderAlert(err.Error(), c)
	}

	dealListType := view.DealListUserDeals
	if len(dealerId) > 0 {
		dealListType = view.DealListDealer
	}

	hideName := c.QueryParam("hide_name") == "true"
	canEdit := c.QueryParam("can_edit") == "true"
	showStatistics := c.QueryParam("show_statistics") == "true"

	actionButton := view.ActionButtonFavoriteToggle
	if canEdit {
		actionButton = view.ActionButtonEdit
	}
	if showStatistics {
		actionButton = view.ActionButtonStatistics
	}

	return view.Render(view.DealsList(dealListType, user, headers, hideName, actionButton), c)
}

type DealJson struct {
	Location model.Point
	Color    string
}

func getDealsAsJson(c echo.Context) error {
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	extent := c.QueryParam("extent")
	boundingBoxFilter := service.BoundingBoxDealFilter{
		BoundingBox: extent,
	}

	deals, err := service.GetActiveDeals(boundingBoxFilter, user)
	if err != nil {
		zap.L().Sugar().Error("can't get deals: ", err)
		return nil
	}

	var dealJson []DealJson
	for _, deal := range deals {
		location, err := model.NewPointFromString(deal.Location)
		if err != nil {
			zap.L().Sugar().Error("can't convert location string to model.Point: ", err)
		}
		jsonEntry := DealJson{
			Location: location,
			Color:    model.GetColorById(deal.CategoryId),
		}
		dealJson = append(dealJson, jsonEntry)
	}

	return c.JSON(http.StatusOK, dealJson)
}

func getDealDetails(c echo.Context) error {
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	dealId := c.Param("id")
	likes := service.GetDealLikes(dealId)
	imageUrls, err := service.GetDealImageUrls(dealId)
	if err != nil {
		zap.L().Sugar().Error("can't get deal image urls: ", err)
		return c.String(http.StatusNotFound, "Konnte Deal Details nicht laden, bitte versuche es später nochmal.")
	}

	details, err := service.GetDealDetails(dealId)
	if err != nil {
		zap.L().Sugar().Error("can't get deal details: ", err)
		return c.String(http.StatusNotFound, "Konnte Deal Details nicht laden, bitte versuche es später nochmal.")
	}

	lang := service.GetLanguageFromCookie(c)

	return view.Render(view.DealDetailsFooter(details, imageUrls, user.IsDealer, strconv.Itoa(likes), lang), c)
}

func toggleDealLike(c echo.Context) error {
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	dealId := c.Param("id")
	doToggle := c.QueryParam("toggle") != "false"
	likes := 0

	if doToggle && user.IsBasicUser {
		return view.RenderInfoTranslated("info.only_for_pro_member", c)
	}

	if doToggle {
		likes = service.ToggleLikes(dealId, user.ID.String())
	} else {
		likes = service.GetDealLikes(dealId)
	}

	isLiked := service.IsDealLiked(dealId, user.ID.String())

	return view.Render(view.Likes(dealId, isLiked, strconv.Itoa(likes)), c)
}

func getReportModal(c echo.Context) error {
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return view.RenderAlertTranslated("alert.can_t_report_deal", c)
	}

	if user.IsBasicUser {
		return view.RenderInfoTranslated("info.only_for_pro_member", c)
	}

	dealId := c.Param("id")
	report, err := service.GetDealReport(dealId, user.ID.String())
	if err != nil {
		zap.L().Sugar().Error("can't get deal report reason: ", err)
	}

	lang := service.GetLanguageFromCookie(c)

	return view.Render(view.DealReportModal(dealId, report.Reason, report.Title, lang), c)
}

func saveReport(c echo.Context) error {
	userId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		zap.L().Sugar().Error("can't save deal report -> no user ID: ", err)
		return view.RenderAlertTranslated("alert.can_t_report_deal", c)
	}

	reason := c.FormValue("reason")
	if len(reason) == 0 {
		zap.L().Sugar().Error("can't save deal report -> no reason")
		return view.RenderAlertTranslated("alert.report_cause", c)
	}

	dealId := c.Param("id")
	err = service.SaveDealReport(dealId, userId.String(), reason)
	if err != nil {
		zap.L().Sugar().Error("can't save deal report: ", err)
		return view.RenderAlertTranslated("alert.can_t_save_report", c)
	}

	return nil
}

func getDealFavoriteButton(c echo.Context) error {
	userId, _ := service.GetUserIdFromCookie(c)
	dealId := c.Param("id")
	isFavorite := service.IsDealFavorite(dealId, userId.String())

	return view.Render(view.DealFavoriteToggleButton(dealId, isFavorite), c)
}

func toggleFavorite(c echo.Context) error {
	user, _ := service.GetUserFromCookie(c)
	dealId := c.Param("id")

	if user.IsBasicUser {
		return view.RenderInfoTranslated("info.only_for_pro_member", c)
	}

	isFavorite := service.ToggleFavorite(dealId, user.ID.String())

	c.Response().Header().Add("HX-Trigger", "updateFavDealsCountBadge")

	return view.Render(view.DealFavoriteToggleButton(dealId, isFavorite), c)
}

func removeFavorite(c echo.Context) error {
	userId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	dealId := c.Param("id")
	err = service.RemoveDealFavorite(dealId, userId.String())
	if err != nil {
		zap.L().Sugar().Error("can't remove deal favorite: ", err)
	}

	return nil
}

func getImageZoomModal(c echo.Context) error {
	dealId := c.Param("dealId")
	startIndex, err := strconv.Atoi(c.QueryParam("index"))
	if err != nil {
		startIndex = 0
	}

	imageUrls, err := service.GetDealImageUrls(dealId)
	if err != nil {
		zap.L().Sugar().Error("can't get deal images: ", err)
		return view.RenderAlertTranslated("alert.can_t_load_deal_images", c)
	}

	return view.Render(view.ImageZoomModal(imageUrls, startIndex), c)
}
