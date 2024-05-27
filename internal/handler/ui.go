package handler

import (
	"github.com/bitstorm-tech/cockaigne/internal/service"
	"github.com/bitstorm-tech/cockaigne/internal/view"
	"github.com/labstack/echo/v4"
)

func RegisterUiHandlers(e *echo.Echo) {
	e.GET("/ui/header", func(c echo.Context) error {
		user, _ := service.GetUserFromCookie(c)

		if len(user.Language) == 0 {
			user.Language = service.LanguageCodeDe
		}

		return view.Render(view.Header(user), c)
	})

	e.GET("/ui/footer", func(c echo.Context) error {
		isAuthenticated := service.IsAuthenticated(c)

		if !isAuthenticated {
			return nil
		}

		isDealer := service.IsDealer(c)

		return view.Render(view.Footer(isDealer), c)
	})

	e.DELETE("/ui/remove", func(c echo.Context) error {
		return nil
	})
}
