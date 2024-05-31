package handler

import (
	"github.com/bitstorm-tech/cockaigne/internal/model"
	"github.com/bitstorm-tech/cockaigne/internal/redirect"
	"github.com/bitstorm-tech/cockaigne/internal/service"
	"github.com/bitstorm-tech/cockaigne/internal/view"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

func RegisterMapHandlers(e *echo.Echo) {
	e.GET("/map", openMap)
	e.GET("/ui/map/filter-modal", openFilterModal)
	e.GET("/ui/map/location-modal", openLocationModal)
}

func openMap(c echo.Context) error {
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	if user.IsDealer {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	if user.IsBasicUser {
		filter := service.GetBasicUserFilter(user.ID.String())
		return view.Render(view.Map(filter.SearchRadiusInMeters, filter.UseLocationService, filter.Location), c)
	}

	acc, err := service.GetAccount(user.ID.String())
	if err != nil {
		zap.L().Sugar().Error("can't get account: ", err)
	}

	location, err := model.NewPointFromString(acc.Location.String)
	if err != nil {
		zap.L().Sugar().Error("can't create new point from database location: ", err)
	}

	return view.Render(view.Map(acc.SearchRadiusInMeters, acc.UseLocationService, location), c)
}

func openFilterModal(c echo.Context) error {
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	categories := service.GetCategories()
	redirectAfterSave := c.QueryParam("redirect-after-save")
	lang := service.GetLanguageFromCookie(c)

	if user.IsBasicUser {
		basicUserFilter := service.GetBasicUserFilter(user.ID.String())
		return view.Render(
			view.FilterModal(categories, basicUserFilter.SelectedCategories, basicUserFilter.SearchRadiusInMeters, redirectAfterSave, lang),
			c,
		)
	}

	acc, err := service.GetAccount(user.ID.String())
	if err != nil {
		zap.L().Sugar().Errorf("can't get account: %v", err)
		return view.RenderAlert("Leider können die Filter gerade nicht geladen werden, bitte versuche es später noch einmal.", c)
	}

	favCategoryIds := service.GetFavoriteCategoryIds(user.ID)

	return view.Render(view.FilterModal(categories, favCategoryIds, acc.SearchRadiusInMeters, redirectAfterSave, lang), c)
}

func openLocationModal(c echo.Context) error {
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	useLocationService := false

	if user.IsBasicUser {
		useLocationService = service.GetBasicUserFilter(user.ID.String()).UseLocationService
	} else {
		acc, err := service.GetAccount(user.ID.String())
		if err != nil {
			zap.L().Sugar().Error("can't get account: ", err)
			return view.RenderAlert("Leider ist uns ein Fehler unterlaufen, bitte versuche es später noch einmal.", c)
		}
		useLocationService = acc.UseLocationService
	}

	lang := service.GetLanguageFromCookie(c)

	return view.Render(view.LocationModal(useLocationService, lang), c)
}
