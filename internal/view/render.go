package view

import (
	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(t templ.Component, c echo.Context) error {
	return t.Render(context.TODO(), c.Response().Writer)
}

func RenderToTarget(t templ.Component, c echo.Context, target string) error {
	c.Response().Header().Add("HX-Retarget", target)
	c.Response().Header().Add("HX-Reswap", "innerHTML")
	return t.Render(context.TODO(), c.Response().Writer)
}
