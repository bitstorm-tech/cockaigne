package user

import "github.com/gofiber/fiber/v2"

func Register(app *fiber.App) {
	app.Get("/user", func(c *fiber.Ctx) error {
		return c.Render("pages/user", nil, "layouts/main")
	})
}
