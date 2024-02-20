package handler

import (
	"github.com/bitstorm-tech/cockaigne/internal/service"
	"github.com/bitstorm-tech/cockaigne/internal/view"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RegisterMapHandlers(e *echo.Echo) {
	e.GET("/map", func(c echo.Context) error {
		userId, _ := service.ParseUserId(c)
		acc, err := service.GetAccount(userId.String())
		if err != nil {
			zap.L().Sugar().Errorf("can't get account: %v", err)
		}

		return view.Render(view.Map(acc.SearchRadiusInMeters, acc.UseLocationService, acc.Location.String), c)
	})

	e.GET("/ui/map/filter-modal", func(c echo.Context) error {
		userId, _ := service.ParseUserId(c)

		acc, err := service.GetAccount(userId.String())

		if err != nil {
			zap.L().Sugar().Errorf("can't get account: %v", err)
		}

		categories := service.GetCategories()
		favCategoryIds := service.GetFavoriteCategoryIds(userId)

		return view.Render(view.FilterModal(categories, favCategoryIds, acc.SearchRadiusInMeters), c)
	})

	e.GET("/ui/map/location-modal", func(c echo.Context) error {
		// userId, err := service.ParseUserId(c)
		// if err != nil {
		// 	return redirect.Login(c)
		// }
		//
		// acc, err := service.GetAccount(userId.String())
		// if err != nil {
		// 	zap.L().Sugar().Error("can't get account: ", err)
		//
		// }
		return view.Render(view.LocationModal(), c)
	})
}
