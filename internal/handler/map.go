package handler

import (
	"github.com/bitstorm-tech/cockaigne/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func RegisterMapHandlers(app *fiber.App) {
	app.Get("/map", func(c *fiber.Ctx) error {
		userId, _ := service.ParseUserId(c)
		acc, err := service.GetAccount(userId.String())
		if err != nil {
			log.Errorf("can't get account: %v", err)
		}

		return c.Render("pages/map", fiber.Map{
			"searchRadius":       acc.SearchRadiusInMeters,
			"useLocationService": acc.UseLocationService,
			"location":           acc.Location.String,
		}, "layouts/main")
	})

	app.Get("/ui/map/filter-modal", func(c *fiber.Ctx) error {
		userId, _ := service.ParseUserId(c)

		acc, err := service.GetAccount(userId.String())

		if err != nil {
			log.Errorf("can't get account: %v", err)
		}

		categories := service.GetCategories()
		favCategoryIds := service.GetFavoriteCategoryIds(userId)

		return c.Render(
			"fragments/map/filter-modal",
			fiber.Map{
				"titel":          "Filter",
				"searchRadius":   acc.SearchRadiusInMeters,
				"categories":     categories,
				"favCategoryIds": favCategoryIds,
			},
		)
	})

	app.Get("/ui/map/location-modal", func(c *fiber.Ctx) error {
		return c.Render("fragments/map/location-modal", fiber.Map{"titel": "Standort festlegen"})
	})
}
