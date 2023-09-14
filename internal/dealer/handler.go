package dealer

import "github.com/gofiber/fiber/v2"

func Register(app *fiber.App) {
	app.Get("/dealer/:dealerId", func(c *fiber.Ctx) error {
		return c.Render("pages/dealer", nil, "layouts/main")
	})
}
