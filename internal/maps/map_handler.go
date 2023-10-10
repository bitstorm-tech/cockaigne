package maps

import (
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	app.Get("/map", func(c *fiber.Ctx) error {
		return c.Render("pages/map", nil, "layouts/main")
	})

	app.Get("/ui/map/filter-modal", func(c *fiber.Ctx) error {
		return c.Render("partials/modal", nil)
	})
}
