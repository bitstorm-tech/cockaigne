package redirect

import (
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	return To("/login", c)
}

func To(location string, c echo.Context) error {
	c.Response().Header().Add("HX-Redirect", location)
	return nil
}
