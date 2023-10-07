package ui

import (
	"github.com/bitstorm-tech/cockaigne/internal/auth"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	app.Get("/ui/partials/header", func(c *fiber.Ctx) error {
		isAuthenticated := auth.IsAuthenticated(c)

		return c.Render("partials/header", fiber.Map{
			"isAuthenticated": isAuthenticated,
		})
	})

	app.Get("/ui/partials/footer", func(c *fiber.Ctx) error {
		isAuthenticated := auth.IsAuthenticated(c)

		if !isAuthenticated {
			return c.Render("", nil)
		}

		isDealer := auth.IsDealer(c)

		return c.Render("partials/footer", fiber.Map{
			"isDealer": isDealer,
		})
	})

	app.Get("/ui/partials/alert", func(c *fiber.Ctx) error {
		return c.Render("partials/alert", fiber.Map{
			"message": "Läuft doch eigentlich ganz gut, oder?",
		})
	})

	app.Delete("/ui/remove", func(c *fiber.Ctx) error {
		return c.SendString("")
	})
}
