package handler

import (
	"github.com/bitstorm-tech/cockaigne/internal/model"
	"github.com/bitstorm-tech/cockaigne/internal/redirect"
	"github.com/bitstorm-tech/cockaigne/internal/service"
	"github.com/bitstorm-tech/cockaigne/internal/view"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RegisterMapHandlers(e *echo.Echo) {
	e.GET("/map", openMap)
	e.GET("/ui/map/filter-modal", openFilterModal)
	e.GET("/ui/map/location-modal", openLocationModal)
}

func openMap(c echo.Context) error {
	user, err := service.ParseUser(c)
	if err != nil {
		return redirect.Login(c)
	}

	if user.IsBasicUser {
		filter := service.GetBasicUserFilter(user.ID)
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
	userId, _ := service.ParseUserId(c)

	acc, err := service.GetAccount(userId.String())

	if err != nil {
		zap.L().Sugar().Errorf("can't get account: %v", err)
	}

	categories := service.GetCategories()
	favCategoryIds := service.GetFavoriteCategoryIds(userId)

	return view.Render(view.FilterModal(categories, favCategoryIds, acc.SearchRadiusInMeters), c)
}

func openLocationModal(c echo.Context) error {
	user, err := service.ParseUser(c)
	if err != nil {
		return redirect.Login(c)
	}

	useLocationService := false

	if user.IsBasicUser {
		useLocationService = service.GetBasicUserFilter(user.ID).UseLocationService
	} else {
		acc, err := service.GetAccount(user.ID.String())
		if err != nil {
			zap.L().Sugar().Error("can't get account: ", err)
			return view.RenderAlert("Leider ist uns ein Fehler unterlaufen, bitte versuche es später noch einmal.", c)
		}
		useLocationService = acc.UseLocationService
	}

	return view.Render(view.LocationModal(useLocationService), c)
}
