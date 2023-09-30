package header

import "github.com/gofiber/fiber/v2"

func Register(app *fiber.App) {
	app.Get("/ui/header", func(c *fiber.Ctx) error {
		showMenu := c.Query("showMenu")
		return c.Render("partials/header", fiber.Map{
			"showMenu": showMenu == "true",
		})
	})

	app.Get("/partials/alert", func(c *fiber.Ctx) error {
		return c.Render("partials/alert", fiber.Map{
			"message": "Läuft doch eigentlich ganz gut, oder?",
		})
	})

	app.Delete("/ui/remove", func(c *fiber.Ctx) error {
		return c.SendString("")
	})
}
