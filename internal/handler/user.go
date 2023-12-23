package handler

import (
	"github.com/bitstorm-tech/cockaigne/internal/redirect"
	"github.com/bitstorm-tech/cockaigne/internal/service"
	"github.com/bitstorm-tech/cockaigne/internal/view"
	"github.com/labstack/echo/v4"
)

func RegisterUserHandlers(e *echo.Echo) {
	e.GET("/user", getUser)
}

func getUser(c echo.Context) error {
	userId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}

	acc, err := service.GetAccount(userId.String())
	if err != nil {
		c.Logger().Errorf("can't get account: %v", err)
	}

	return view.Render(view.User(acc.ID.String(), acc.Username, "Josef-Frankl-Str.", "31A", "80995", "München", "12"), c)
}
