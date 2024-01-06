package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bitstorm-tech/cockaigne/internal/model"
	"github.com/bitstorm-tech/cockaigne/internal/redirect"
	"github.com/bitstorm-tech/cockaigne/internal/service"
	"github.com/bitstorm-tech/cockaigne/internal/view"
	"github.com/labstack/echo/v4"
)

func RegisterDealHandlers(e *echo.Echo) {
	e.GET("/deal/:dealId", openDealCreatePage)
	e.GET("/ui/category-select", getCategorySelect)
	e.GET("/deals/:state", getDealList)
	e.GET("/deal-list/:state", getDealsByState)
	e.GET("/api/deals", getDealsAsJson)
	e.GET("/deal-details/:id", getDealDetails)
	e.GET("/deal-likes/:id", toggleDealLike)
	e.GET("/ui/deals/report-modal/:id", getReportModal)
	e.GET("/deal-favorite-button/:id", getDealFavoriteButton)
	e.GET("/deal-favorite-toggle/:id", toggleFavorite)
	e.GET("/deal-favorites-list", getFavoriteDeals)
	e.GET("/deal-image-zoom-modal/:dealId", getImageZoomModal)
	e.GET("/dealer-favorites-list", getFavoriteDealerDeals)
	e.POST("/deal-report/:id", saveReport)
	e.POST("/deals", saveDeal)
	e.DELETE("/deal-favorite-remove/:id", removeFavorite)
}

func getFavoriteDealerDeals(c echo.Context) error {
	userId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}

	headers, err := service.GetFavoriteDealerDealHeaders(userId.String())
	if err != nil {
		c.Logger().Errorf("can't get favorite dealer deals: %v", err)
		return view.RenderAlert("Kann favorisierte Dealer Deals nicht laden, bitte später nochmal versuchen.", c)
	}

	return view.Render(view.DealsList(headers, false, false, false, false), c)
}

func getFavoriteDeals(c echo.Context) error {
	userId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}

	headers, err := service.GetFavoriteDealHeaders(userId.String())
	if err != nil {
		c.Logger().Errorf("can't get favorite deal headers: %v", err)
		return view.RenderAlert("Kann favorisierte Deals aktuell nicht laden, bitte später nochmal versuchen.", c)
	}

	return view.Render(view.DealsList(headers, false, false, true, true), c)
}

func getDealsByState(c echo.Context) error {
	user, err := service.ParseUser(c)
	if err != nil {
		c.Logger().Errorf("can't parse user: %v", err)
		return redirect.Login(c)
	}

	state := model.ToState(c.Param("state"))
	userIdString := user.ID.String()
	userId := &userIdString

	if !user.IsDealer {
		userId = nil
	}

	headers, err := service.GetDealHeaders(state, userId)
	if err != nil {
		c.Logger().Error(err)
		return view.RenderAlert(err.Error(), c)
	}

	onDealerPage := strings.Contains(c.Request().URL.Path, "dealer")

	return view.Render(view.DealsList(headers, user.IsDealer, onDealerPage, true, false), c)
}

func openDealCreatePage(c echo.Context) error {
	userId, err := service.ParseUserId(c)
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
			return view.RenderAlert("Der Deal konnte leider nicht gefunden werden, bitte versuche es später nochmal.", c)
		}
	}

	return view.Render(view.DealEdit(deal), c)
}

func getCategorySelect(c echo.Context) error {
	categories := service.GetCategories()
	name := c.QueryParam("name")
	selectedParam := c.QueryParam("selected")

	selected := -1
	var err error
	if len(selectedParam) > 0 {
		selected, err = strconv.Atoi(c.QueryParam("selected"))
		if err != nil {
			c.Logger().Errorf("can't parse selected category: %v", err)
			selected = -1
		}
	}

	return view.Render(view.CategorySelect(name, categories, selected), c)
}

func saveDeal(c echo.Context) error {
	userId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}

	deal, errorMessage := model.DealFromRequest(c)
	if len(errorMessage) > 0 {
		return view.RenderAlert(errorMessage, c)
	}

	deal.DealerId = userId
	c.Logger().Debugf("Create deal: %+v", deal)

	dealId, err := service.SaveDeal(deal)
	if err != nil {
		c.Logger().Errorf("can't save deal: %v", err)
		return view.RenderAlert("Leider ist beim Erstellen etwas schief gegangen, bitte versuche es später nochmal.", c)
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.Logger().Errorf("can't get multipart form: %v", err)
		return view.RenderAlert("Leider ist beim Erstellen etwas schief gegangen, bitte versuche es später nochmal.", c)
	}

	for index, file := range form.File["images"] {
		err = service.UploadDealImage(file, dealId.String(), fmt.Sprintf("%d-", index))
		if err != nil {
			c.Logger().Errorf("can't upload deal image: %v", err)
			return view.RenderAlert("Leider ist beim Erstellen etwas schief gegangen, bitte versuche es später nochmal.", c)
		}
	}

	c.Response().Header().Set("HX-Redirect", "/")

	return nil
}

func getDealList(c echo.Context) error {
	user, err := service.ParseUser(c)
	if err != nil {
		c.Logger().Errorf("can't parse user: %v", err)
		return redirect.Login(c)
	}

	state := model.ToState(c.Param("state"))
	userIdString := user.ID.String()
	userId := &userIdString

	if !user.IsDealer {
		userId = nil
	}

	headers, err := service.GetDealHeaders(state, userId)
	if err != nil {
		c.Logger().Error(err)
		return view.RenderAlert(err.Error(), c)
	}

	onDealerPage := strings.Contains(c.Request().URL.Path, "dealer")

	return view.Render(view.DealsList(headers, user.IsDealer, onDealerPage, true, false), c)
}

func getDealsAsJson(c echo.Context) error {
	// extent := c.Query("extent")
	deals, err := service.GetDealsFromView(model.Active, nil)
	if err != nil {
		c.Logger().Errorf("can't get deals: %v", err)
		return nil
	}

	return c.JSON(http.StatusOK, deals)
}

func getDealDetails(c echo.Context) error {
	dealId := c.Param("id")
	likes := service.GetDealLikes(dealId)
	imageUrls, err := service.GetDealImageUrls(dealId)
	if err != nil {
		c.Logger().Errorf("can't get deal image urls: %v", err)
		return c.String(http.StatusNotFound, "Konnte Deal Details nicht laden, bitte versuche es später nochmal.")
	}

	details, err := service.GetDealDetails(dealId)
	if err != nil {
		c.Logger().Errorf("can't get deal details: %v", err)
		return c.String(http.StatusNotFound, "Konnte Deal Details nicht laden, bitte versuche es später nochmal.")
	}

	return view.Render(view.DealDetailsFooter(details, imageUrls, true, strconv.Itoa(likes)), c)
}

func toggleDealLike(c echo.Context) error {
	userId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}
	dealId := c.Param("id")
	doToggle := c.QueryParam("toggle") != "false"
	likes := 0
	if doToggle {
		likes = service.ToggleLikes(dealId, userId.String())
	} else {
		likes = service.GetDealLikes(dealId)
	}

	isLiked := service.IsDealLiked(dealId, userId.String())

	return view.Render(view.Likes(dealId, isLiked, strconv.Itoa(likes)), c)
}

func getReportModal(c echo.Context) error {
	dealId := c.Param("id")
	reporterId, err := service.ParseUserId(c)
	if err != nil {
		return view.RenderAlert("Nur angemeldete User können einen Deal melden", c)
	}

	report, err := service.GetDealReport(dealId, reporterId.String())
	if err != nil {
		c.Logger().Errorf("can't get deal report reason: %v", err)
	}

	return view.Render(view.DealReportModal(dealId, report.Reason, report.Title), c)
}

func saveReport(c echo.Context) error {
	userId, err := service.ParseUserId(c)
	if err != nil {
		c.Logger().Error("can't save deal report -> no user ID")
		return view.RenderAlert("Nur angemeldete User können einen Deal melden", c)
	}

	reason := c.FormValue("reason")
	if len(reason) == 0 {
		c.Logger().Error("can't save deal report -> no reason")
		return view.RenderAlert("Bitte gib an, was an dem Deal nicht passt", c)
	}

	dealId := c.Param("id")
	err = service.SaveDealReport(dealId, userId.String(), reason)
	if err != nil {
		c.Logger().Errorf("can't save deal report: %v", err)
		return view.RenderAlert("Deal konnte leider nicht gemeldet werden, bitte versuche es später noch einmal.", c)
	}

	return nil
}

func getDealFavoriteButton(c echo.Context) error {
	userId, _ := service.ParseUserId(c)
	dealId := c.Param("id")
	isFavorite := service.IsDealFavorite(dealId, userId.String())

	return view.Render(view.DealFavoriteToggleButton(dealId, isFavorite), c)
}

func toggleFavorite(c echo.Context) error {
	userId, _ := service.ParseUserId(c)
	dealId := c.Param("id")
	isFavorite := service.ToggleFavorite(dealId, userId.String())

	return view.Render(view.DealFavoriteToggleButton(dealId, isFavorite), c)
}

func removeFavorite(c echo.Context) error {
	userId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}

	dealId := c.Param("id")
	err = service.RemoveDealFavorite(dealId, userId.String())
	if err != nil {
		c.Logger().Errorf("can't remove deal favorite: %v", err)
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
		c.Logger().Errorf("can't get deal images: %v", err)
		return view.RenderAlert("Kann Deal Bilder momentan nicht laden, bitte versuche es später nochmal.", c)
	}

	return view.Render(view.ImageZoomModal(imageUrls, startIndex), c)
}
