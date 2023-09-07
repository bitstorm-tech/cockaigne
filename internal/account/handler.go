package account

import (
	"log"

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
		return c.Render("account/index", nil, "layouts/main")
	})

	app.Get("/create-account", func(c *fiber.Ctx) error {
		return c.Render("account/create-account", nil, "layouts/main")
	})

	app.Post("/create-account", handleCreateAccount)
}

func handleCreateAccount(c *fiber.Ctx) error {
	newAccount := Account{}

	err := c.BodyParser(&newAccount)
	if err != nil {
		log.Println("Can't parse request body:", err)
		return c.SendString("ERROR")
	}

	log.Printf("Create new account: %v\n", newAccount)

	return c.SendString("DONE")
}
