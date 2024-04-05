package service

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func GetBaseUrl(c echo.Context) string {
	return fmt.Sprintf("%s://%s", c.Scheme(), c.Request().Host)
}
