package dealer

import (
	"github.com/bitstorm-tech/cockaigne/internal/account"
	"github.com/bitstorm-tech/cockaigne/internal/deal"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Register(app *fiber.App) {
	app.Get("/dealer/:dealerId", func(c *fiber.Ctx) error {
		dealerId := c.Params("dealerId")
		acc, err := account.GetAccount(dealerId)

		if err != nil {
			log.Errorf("can't find dealer (%s): %v", dealerId, err)
			return c.Status(fiber.StatusNotFound).SendString("Not Found")
		}

		category, err := deal.GetCategory(int(acc.DefaultCategory.Int32))

		if err != nil {
			log.Errorf("can't find category (id=%s): %v", acc.DefaultCategory.Int32, err)
			return c.Status(fiber.StatusNotFound).SendString("Not Found")
		}

		return c.Render("pages/dealer", fiber.Map{"account": acc, "category": category.Name}, "layouts/main")
	})
}
