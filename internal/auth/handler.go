package auth

import (
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Register(app *fiber.App) {
	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("pages/login", nil)
	})

	app.Get("/signup", func(c *fiber.Ctx) error {
		return c.Render("pages/signup", nil, "layouts/annonym")
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

	log.Debugf("New account: %+v", acc.Email)

	err = persistence.DB.Create(&acc).Error

	if err != nil {
		return c.Render("partials/alert", fiber.Map{"message": err.Error()})
	}

	return c.Redirect("/login")
}
