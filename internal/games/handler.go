package games

import (
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GameMetadata struct {
	gorm.Model
	ID              string
	Name            string
	Description     string
	ShowDescription bool
}

func (GameMetadata) TableName() string {
	return "games_metadata"
}

func Register(app *fiber.App) {
	app.Get("/games", func(c *fiber.Ctx) error {
		return c.Render("pages/games", nil, "layouts/main")
	})

	app.Get("/game-list", func(c *fiber.Ctx) error {
		clicked := c.Query("clickedGame")

		var gamesMetadata []GameMetadata
		persistence.DB.Find(&gamesMetadata)

		for gameListItem := range gamesMetadata {
			if gamesMetadata[gameListItem].Name == clicked {
				gamesMetadata[gameListItem].ShowDescription = !gamesMetadata[gameListItem].ShowDescription
			} else {
				gamesMetadata[gameListItem].ShowDescription = false
			}
		}

		return c.Render("pages/game-list", fiber.Map{"games": gamesMetadata, "count": len(gamesMetadata)})
	})
}
