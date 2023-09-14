package dealer

import (
	"github.com/bitstorm-tech/cockaigne/internal/account"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	app.Get("/dealer/:dealerId", func(c *fiber.Ctx) error {
		acc := account.Account{}
		err := persistence.DB.First(&acc, "id = ?", c.Params("dealerId")).Error

		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Not Found")
		}

		return c.Render("pages/dealer", fiber.Map{"account": acc}, "layouts/main")
	})
}
