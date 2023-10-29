package deal

import (
	"github.com/bitstorm-tech/cockaigne/internal/account"
	"github.com/bitstorm-tech/cockaigne/internal/ui"
	"strings"

	"github.com/bitstorm-tech/cockaigne/internal/auth/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Register(app *fiber.App) {
	app.Get("/deal/:dealId", func(c *fiber.Ctx) error {
		userId, err := jwt.ParseUserId(c)
		if err != nil {
			return c.Redirect("/login")
		}

		dealId := c.Params("dealId")

		var deal Deal
		if strings.EqualFold(dealId, "new") {
			deal = NewDeal()
			deal.CategoryId = account.GetDefaultCategoryId(userId)
		} else {
			var err error
			deal, err = GetDeal(dealId)
			if err != nil {
				return c.Status(fiber.StatusNotFound).SendString("Not Found")
			}
		}

		return c.Render("pages/edit-deal", fiber.Map{"deal": deal}, "layouts/main")
	})

	app.Get("/ui/category-select", func(c *fiber.Ctx) error {
		categories := GetCategories()
		name := c.Query("name", "Kategorie")
		selected := c.Query("selected", "-1")

		return c.Render("partials/category-select", fiber.Map{"categories": categories, "name": name, "selected": selected})
	})

	app.Post("/api/deals", func(c *fiber.Ctx) error {
		userId, err := jwt.ParseUserId(c)
		if err != nil {
			return c.Redirect("/login")
		}

		deal, errorMessage := NewDealFromRequest(c)
		if len(errorMessage) > 0 {
			return ui.ShowAlert(c, errorMessage)
		}

		deal.DealerId = userId
		log.Debugf("Create deal: %+v", deal)

		if err := SaveDeal(deal); err != nil {
			return ui.ShowAlert(c, err.Error())
		}

		c.Set("HX-Redirect", "/")

		return nil
	})

	app.Get("/deals/:state", func(c *fiber.Ctx) error {
		userId, err := jwt.ParseUserId(c)
		if err != nil {
			return c.Redirect("/login")
		}

		state := ToState(c.Params("state", "active"))
		userIdString := userId.String()

		deals, err := GetDealsFromView(state, &userIdString)
		if err != nil {
			log.Error(err)
			return ui.ShowAlert(c, err.Error())
		}

		return c.Render("partials/deals-list", fiber.Map{"deals": deals})
	})

	app.Get("/api/deals", func(c *fiber.Ctx) error {
		// extent := c.Query("extent")
		deals, err := GetDealsFromView(Active, nil)
		if err != nil {
			log.Errorf("can't get deals: %v", err)
			return nil
		}

		return c.JSON(deals)
	})
}
