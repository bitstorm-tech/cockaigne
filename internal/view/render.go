package view

import (
	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(t templ.Component, c echo.Context) error {
	return t.Render(context.TODO(), c.Response().Writer)
}
