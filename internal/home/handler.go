package home

import (
	"github.com/bitstorm-tech/cockaigne/internal/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Register(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		jwt := c.Cookies("jwt")
		claims, err := auth.ParseJwtToken(jwt)

		if err != nil {
			log.Errorf("Can't parse JWT token: %+v", err)
			return c.Redirect("/login")
		}

		if claims["isDealer"] == true {
			id := claims["sub"].(string)
			return c.Redirect("/dealer/" + id)
		}

		return c.Redirect("/user")
	})
}
