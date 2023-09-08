package account

import (
	"github.com/gofiber/fiber/v2"
)

type Account struct {
	ID       string
	Email    string
	Nickname string
	Password string
}

func Register(app *fiber.App) {
	app.Get("/account", func(c *fiber.Ctx) error {
		return c.Render("pages/account", nil, "layouts/main")
	})
}
