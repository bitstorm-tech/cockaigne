package handler

import (
	"fmt"

	"github.com/bitstorm-tech/cockaigne/internal/redirect"
	"github.com/bitstorm-tech/cockaigne/internal/service"
	"github.com/bitstorm-tech/cockaigne/internal/view"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RegisterUserHandlers(e *echo.Echo) {
	e.GET("/user", getUser)
	e.GET("/deals-count-badge", getDealsCountBadge)
	e.GET("/favorite-deals-count-badge", getFavoriteDealsCountBadge)
	e.GET("/favorite-dealers-count-badge", getFavoriteDealersCountBadge)
}

func getDealsCountBadge(c echo.Context) error {
	return view.Render(view.DealsCountBadge("1"), c)
}

func getFavoriteDealersCountBadge(c echo.Context) error {
	userId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}

	headers, err := service.GetFavoriteDealerDealHeaders(userId.String())
	if err != nil {
		zap.L().Sugar().Error("can't get favorite dealer deal headers: ", err)
	}

	count := fmt.Sprintf("%d", len(headers))

	return view.Render(view.FavoriteDealerCountBadge(count), c)
}

func getFavoriteDealsCountBadge(c echo.Context) error {
	userId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}

	count, err := service.GetFavoriteDealsCount(userId.String())
	if err != nil {
		zap.L().Sugar().Error("can't get favorite deals count: ", err)
	}

	countString := fmt.Sprintf("%d", count)

	return view.Render(view.FavoriteDealsCountBadge(countString), c)
}

func getUser(c echo.Context) error {
	userId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}

	acc, err := service.GetAccount(userId.String())
	if err != nil {
		zap.L().Sugar().Error("can't get account: ", err)
	}

	params := view.UserHeaderParameters{
		ID:          acc.ID.String(),
		Username:    acc.Username,
		Street:      "Suche aktuelle Position ...",
		HouseNumber: "",
		Zip:         "",
		City:        "",
	}

	return view.Render(view.User(params, acc.UseLocationService, acc.Location.String), c)
}
