package header

import "github.com/gofiber/fiber/v2"

func Register(app *fiber.App) {
	app.Get("/header", func(c *fiber.Ctx) error {
		showMenu := c.Query("showMenu")
		return c.Render("header/header", fiber.Map{
			"showMenu": showMenu == "true",
		})
	})
}
