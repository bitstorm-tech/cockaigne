package maps

import (
	"github.com/bitstorm-tech/cockaigne/internal/account"
	"github.com/bitstorm-tech/cockaigne/internal/auth/jwt"
	"github.com/bitstorm-tech/cockaigne/internal/deal"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Register(app *fiber.App) {
	app.Get("/map", func(c *fiber.Ctx) error {
		userId, _ := jwt.ParseUserId(c)

		searchRadius := account.GetSearchRadius(userId)

		return c.Render("pages/map", fiber.Map{
			"searchRadius": searchRadius,
		}, "layouts/main")
	})

	app.Get("/ui/map/filter-modal", func(c *fiber.Ctx) error {
		userId, _ := jwt.ParseUserId(c)

		acc, err := account.GetAccount(userId.String())

		if err != nil {
			log.Errorf("can't get account: %v", err)
		}

		categories := deal.GetCategories()
		favCategoryIds := account.GetFavoriteCategoryIds(userId)

		return c.Render("partials/map/filter-modal", fiber.Map{
			"titel":          "Filter",
			"searchRadius":   acc.SearchRadiusInMeters,
			"categories":     categories,
			"favCategoryIds": favCategoryIds,
		}, "layouts/modal")
	})
}
