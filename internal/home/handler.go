package home

import "github.com/gofiber/fiber/v2"

func Register(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("user/index", nil, "layouts/main")
	})
}
