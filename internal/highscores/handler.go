package highscores

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

var score = 0

func Register(app *fiber.App) {
	app.Get("/highscores", func(c *fiber.Ctx) error {
		return c.Render("pages/highscores", nil, "layouts/main")
	})

	app.Get("/highscores/increment", func(c *fiber.Ctx) error {
		score++
		return c.SendString(strconv.Itoa(score))
	})

	app.Get("/highscores/decrement", func(c *fiber.Ctx) error {
		score--
		return c.SendString(strconv.Itoa(score))
	})

	app.Get("/highscores/current", func(c *fiber.Ctx) error {
		return c.SendString(strconv.Itoa(score))
	})

	app.Get("/highscores/reset", func(c *fiber.Ctx) error {
		score = 0
		return c.SendString(strconv.Itoa(score))
	})

	app.Get("/highscores/wait", func(c *fiber.Ctx) error {
		time.Sleep(1 * time.Second)
		score = 1337
		return c.SendString(strconv.Itoa(score))
	})
}
