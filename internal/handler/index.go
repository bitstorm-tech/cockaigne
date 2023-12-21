package handler

import (
	"net/http"

	"github.com/bitstorm-tech/cockaigne/internal/service"
	"github.com/labstack/echo/v4"
)

func RegisterIndexHandlers(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		user, err := service.ParseUserEcho(c)

		if err != nil {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}

		if user.IsDealer {
			return c.Redirect(http.StatusTemporaryRedirect, "/dealer/"+user.ID.String())
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/user")
	})
}
