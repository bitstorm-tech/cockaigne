package redirect

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func To(location string, c echo.Context) error {
	c.Response().Header().Add("HX-Redirect", location)
	return nil
}
