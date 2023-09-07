package games

import (
	"github.com/bitstorm-tech/cockaigne/internal/db"
	"github.com/gofiber/fiber/v2"
)

type GameListItem struct {
	ID              string
	Name            string
	Description     string
	ShowDescription bool
}

func Register(app *fiber.App) {
	app.Get("/games", func(c *fiber.Ctx) error {
		return c.Render("games/index", nil, "layouts/main")
	})

	app.Get("/game-list", func(c *fiber.Ctx) error {
		clicked := c.Query("clickedGame")

		rows, err := db.Conn.Query("select id, name, description from games_metadata")
		if err != nil {
			return c.SendString(err.Error())
		}

		var gameListItems []GameListItem

		for rows.Next() {
			var gameListItem GameListItem
			err = rows.Scan(&gameListItem.ID, &gameListItem.Name, &gameListItem.Description)

			if err != nil {
				return c.SendString(err.Error())
			}

			gameListItems = append(gameListItems, gameListItem)
		}

		for gameListItem := range gameListItems {
			if gameListItems[gameListItem].Name == clicked {
				gameListItems[gameListItem].ShowDescription = !gameListItems[gameListItem].ShowDescription
			} else {
				gameListItems[gameListItem].ShowDescription = false
			}
		}

		return c.Render("games/game-list", fiber.Map{"games": gameListItems, "count": len(gameListItems)})
	})
}
