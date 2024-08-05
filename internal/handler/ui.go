package handler

import (
	"strings"

	"github.com/bitstorm-tech/cockaigne/internal/service"
	"github.com/bitstorm-tech/cockaigne/internal/view"
	"github.com/labstack/echo/v4"
)

func RegisterUiHandlers(e *echo.Echo) {
	e.GET("/ui/header", func(c echo.Context) error {
		user, _ := service.GetUserFromCookie(c)
		lang := service.GetLanguageFromCookie(c)

		return view.Render(view.Header(user, lang), c)
	})

	e.GET("/ui/footer", func(c echo.Context) error {
		isAuthenticated := service.IsAuthenticated(c)

		if !isAuthenticated {
			return nil
		}

		isDealer := service.IsDealer(c)

		url := strings.ToLower(c.Request().Header.Get("Referer"))
		activeUrl := view.ActiveUrlHome
		if strings.Contains(url, "map") {
			activeUrl = view.ActiveUrlMap
		} else if strings.Contains(url, "top") {
			activeUrl = view.ActiveUrlTopDeals
		} else if strings.Contains(url, "overview") {
			activeUrl = view.ActiveUrlDealsOverview
		}

		return view.Render(view.Footer(isDealer, activeUrl), c)
	})

	e.DELETE("/ui/remove", func(c echo.Context) error {
		return nil
	})
}
