package ui

import "github.com/gofiber/fiber/v2"

func ShowAlert(c *fiber.Ctx, message string) error {
	return c.Render("fragments/alert", fiber.Map{"message": message})
}
