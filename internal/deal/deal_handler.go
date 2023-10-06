package deal

import (
	"strings"

	"github.com/bitstorm-tech/cockaigne/internal/auth"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Register(app *fiber.App) {
	app.Get("/deal/:dealId", func(c *fiber.Ctx) error {
		dealId := c.Params("dealId")

		var deal Deal
		if strings.EqualFold(dealId, "new") {
			deal = NewDeal()
		} else {
			err := persistence.DB.First(&deal, "id = ?", dealId).Error
			if err != nil {
				return c.Status(fiber.StatusNotFound).SendString("Not Found")
			}
		}

		return c.Render("pages/edit-deal", fiber.Map{"deal": deal}, "layouts/main")
	})

	app.Get("/ui/category-select", func(c *fiber.Ctx) error {
		categories := []Category{}
		err := persistence.DB.Find(&categories).Where("active = true").Error
		if err != nil {
			return c.Render("partials/alert", fiber.Map{"message": err.Error()})
		}
		name := c.Query("name", "Kategorie")

		return c.Render("partials/category-select", fiber.Map{"categories": categories, "name": name})
	})

	app.Post("/api/deals", func(c *fiber.Ctx) error {
		userId, err := auth.ParseUserId(c)
		if err != nil {
			return c.Redirect("/login")
		}

		deal, errorMessage := NewDealFromRequest(c)
		if len(errorMessage) > 0 {
			return c.Render("partials/alert", fiber.Map{"message": errorMessage})
		}

		deal.DealerId = userId
		log.Debugf("Create deal: %+v", deal)
		log.Debugf("Deal start: %+v", deal.Start)

		err = persistence.DB.Save(&deal).Error
		if err != nil {
			return c.Render("partials/alert", fiber.Map{"message": err.Error()})
		}

		c.Set("HX-Redirect", "/")

		return nil
	})
}
