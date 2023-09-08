package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Register(app *fiber.App) {
	app.Get("/auth/login", func(c *fiber.Ctx) error {
		return c.Render("auth/login/index", nil)
	})

	app.Get("/auth/signup", func(c *fiber.Ctx) error {
		return c.Render("auth/signup/index", nil, "layouts/annonym")
	})

	app.Post("/api/signup", signup)
}

func signup(c *fiber.Ctx) error {
	acc := Account{}
	err := c.BodyParser(&acc)

	if err != nil {
		log.Errorf("Error while signeup %v", err)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	log.Debugf("Account: %+v", acc)

	return c.JSON(fiber.Map{
		"message": "signup",
	})
}
