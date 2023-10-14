package maps

import (
	"github.com/bitstorm-tech/cockaigne/internal/account"
	"github.com/bitstorm-tech/cockaigne/internal/auth"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	app.Get("/map", func(c *fiber.Ctx) error {
		userId, _ := auth.ParseUserId(c)

		searchRadius := account.GetSearchRadius(userId)

		return c.Render("pages/map", fiber.Map{
			"searchRadius": searchRadius,
		}, "layouts/main")
	})

	app.Get("/ui/map/filter-modal", func(c *fiber.Ctx) error {
		userId, _ := auth.ParseUserId(c)

		searchRadius := account.GetSearchRadius(userId)

		return c.Render("partials/modal", fiber.Map{
			"searchRadius": searchRadius,
		})
	})
}
