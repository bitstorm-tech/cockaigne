package handler

import (
	"fmt"
	"strings"

	"github.com/bitstorm-tech/cockaigne/internal/model"
	"github.com/bitstorm-tech/cockaigne/internal/redirect"
	"github.com/bitstorm-tech/cockaigne/internal/service"
	"github.com/bitstorm-tech/cockaigne/internal/view"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RegisterUserHandlers(e *echo.Echo) {
	e.GET("/user", getUser)
	e.GET("/favorite-deals-count-badge", getFavoriteDealsCountBadge)
	e.GET("/favorite-dealers-count-badge", getFavoriteDealersCountBadge)
	e.POST("/show-new-deals-button", showNewDealsButton)
}

func showNewDealsButton(c echo.Context) error {
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	dealIdsString := c.FormValue("dealIds")
	var dealIds []string

	if len(dealIdsString) > 0 {
		dealIds = strings.Split(dealIdsString, ",")
	}

	newDealsAvailable, err := service.NewDealsAvailable(user, dealIds)
	if err != nil {
		zap.L().Sugar().Error("can't check if new deals are available: ", err)
		return nil
	}

	if newDealsAvailable {
		return view.Render(view.NewDealsButton(), c)
	}

	return nil
}

func getFavoriteDealersCountBadge(c echo.Context) error {
	userId, err := service.GetUserIdFromCookie(c)
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
	userId, err := service.GetUserIdFromCookie(c)
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
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	if user.IsBasicUser {
		filter := service.GetBasicUserFilter(user.ID.String())
		return view.Render(view.User(user.ID.String(), "Basic", false, true, filter.Location), c)
	}

	acc, err := service.GetAccount(user.ID.String())
	if err != nil {
		zap.L().Sugar().Error("can't get account: ", err)
	}

	location, err := model.NewPointFromString(acc.Location.String)
	if err != nil {
		zap.L().Sugar().Errorf("can't create new point from account location (%s): %v", acc.Location.String, err)
	}

	return view.Render(view.User(acc.ID.String(), acc.Username, acc.UseLocationService, false, location), c)
}
