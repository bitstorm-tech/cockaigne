package deal

import (
	"strings"

	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/gofiber/fiber/v2"
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
}
