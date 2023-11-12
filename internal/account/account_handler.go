package account

import (
	"github.com/bitstorm-tech/cockaigne/internal/auth/jwt"
	"github.com/bitstorm-tech/cockaigne/internal/ui"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Register(app *fiber.App) {
	app.Get("/account", func(c *fiber.Ctx) error {
		return c.Render("pages/account", nil, "layouts/main")
	})

	app.Post("/api/accounts/filter", updateFilter)
	app.Post("/api/accounts/use-location-service", updateUseLocationService)
}

func updateFilter(c *fiber.Ctx) error {
	userId, _ := jwt.ParseUserId(c)

	updateFilterRequest := UpdateFilterRequest{}

	err := c.BodyParser(&updateFilterRequest)
	if err != nil {
		log.Errorf("can't parse filter update request: %v", err)
	}

	if err := UpdateSearchRadius(userId, updateFilterRequest.SearchRadiusInMeters); err != nil {
		log.Errorf("can't update accounts search_radius_in_meters: %v", err)
		return ui.ShowAlert(c, "Fehler beim Verarbeiten der Filteränderung")
	}

	if err := UpdateSelectedCategories(userId, updateFilterRequest.FavoriteCategoryIds); err != nil {
		log.Errorf("can't update selected categories: %s", err)
		return ui.ShowAlert(c, "Fehler beim Verarbeiten der Filteränderung")
	}

	return nil
}

func updateUseLocationService(c *fiber.Ctx) error {
	userId, err := jwt.ParseUserId(c)
	if err != nil {
		return c.Redirect("/login")
	}

	useLocationService := c.FormValue("use-location-service")
	err = UpdateUseLocationService(userId.String(), useLocationService == "on")
	if err != nil {
		log.Errorf("can't save use location service: %v", err)
		return ui.ShowAlert(c, "Kann Einstellung leider nicht speichern. Bitte später nochmal versuchen.")
	}

	return nil
}
