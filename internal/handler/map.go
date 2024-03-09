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
	userId, _ := service.ParseUserId(c)
	acc, err := service.GetAccount(userId.String())
	if err != nil {
		zap.L().Sugar().Errorf("can't get account: %v", err)
		return nil
	}

	location, err := model.NewPointFromString(acc.Location.String)
	if err != nil {
		zap.L().Sugar().Error("can't create new point from database location: ", err)
		return nil
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
	userId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}

	acc, err := service.GetAccount(userId.String())
	if err != nil {
		zap.L().Sugar().Error("can't get account: ", err)
		return view.RenderAlert("Momentan", c)
	}

	return view.Render(view.LocationModal(acc.UseLocationService), c)
}
