package redirect

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}
