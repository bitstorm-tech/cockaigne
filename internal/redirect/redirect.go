package redirect

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Login(c echo.Context) error {
	return To("/login", c)
}

func To(location string, c echo.Context) error {
	return c.Redirect(http.StatusTemporaryRedirect, location)
}
