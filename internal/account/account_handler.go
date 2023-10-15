package account

import (
	"github.com/bitstorm-tech/cockaigne/internal/auth/jwt"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/bitstorm-tech/cockaigne/internal/ui"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Register(app *fiber.App) {
	app.Get("/account", func(c *fiber.Ctx) error {
		return c.Render("pages/account", nil, "layouts/main")
	})

	app.Post("/api/accounts/filter", updateFilter)
}

func updateFilter(c *fiber.Ctx) error {
	userId, _ := jwt.ParseUserId(c)

	updateFilterRequest := UpdateFilterRequest{}

	err := c.BodyParser(&updateFilterRequest)
	if err != nil {
		log.Errorf("can't parse filter update request: %v", err)
		return ui.ShowAlert(c, "Fehler beim Verarbeiten der Filteränderung")
	}

	persistence.DB.Model(&Account{}).Where("id = ?", userId).Update("search_radius_in_meters", updateFilterRequest.SearchRadiusInMeters)
	persistence.DB.Where("account_id = ?", userId).Delete(&FavoriteCategory{})
	for _, categoryId := range updateFilterRequest.FavoriteCategoryIds {
		persistence.DB.Create(&FavoriteCategory{AccountId: userId, CategoryId: categoryId})
	}

	return nil
}
