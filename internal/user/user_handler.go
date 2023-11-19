package user

import (
	"github.com/bitstorm-tech/cockaigne/internal/account"
	"github.com/bitstorm-tech/cockaigne/internal/auth/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Register(app *fiber.App) {
	app.Get("/user", func(c *fiber.Ctx) error {
		userId, err := jwt.ParseUserId(c)
		if err != nil {
			return c.Redirect("/login")
		}

		acc, err := account.GetAccount(userId.String())
		if err != nil {
			log.Errorf("can't get account: %v", err)
		}

		return c.Render("pages/user", fiber.Map{
			"username":      acc.Username,
			"street":        "Josef-Frankl-Str.",
			"houseNumber":   "31a",
			"zip":           "80995",
			"city":          "München",
			"numberOfDeals": 12,
		}, "layouts/main")
	})
}
