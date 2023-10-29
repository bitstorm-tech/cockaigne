package dealer

import (
	"github.com/bitstorm-tech/cockaigne/internal/account"
	"github.com/bitstorm-tech/cockaigne/internal/auth/jwt"
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
			log.Error(err)
			return c.Status(fiber.StatusNotFound).SendString("Not Found")
		}

		return c.Render("pages/dealer", fiber.Map{"account": acc, "category": category.Name}, "layouts/main")
	})

	app.Get("/deal-overview", func(c *fiber.Ctx) error {
		return c.Render("pages/deal-overview", nil, "layouts/main")
	})

	app.Get("/templates", func(c *fiber.Ctx) error {
		userId, err := jwt.ParseUserId(c)
		if err != nil {
			return c.Redirect("/login")
		}

		templates, err := deal.GetTemplates(userId.String())
		if err != nil {
			log.Errorf("can't get templates for dealer: %s", userId.String())
		}

		return c.Render("pages/templates", fiber.Map{"templates": templates}, "layouts/main")
	})
}
