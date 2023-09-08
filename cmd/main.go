package main

import (
	"log"
	"os"

	"github.com/bitstorm-tech/cockaigne/internal/account"
	"github.com/bitstorm-tech/cockaigne/internal/auth"
	"github.com/bitstorm-tech/cockaigne/internal/games"
	"github.com/bitstorm-tech/cockaigne/internal/header"
	"github.com/bitstorm-tech/cockaigne/internal/highscores"
	"github.com/bitstorm-tech/cockaigne/internal/home"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	hostAndPort := host + ":" + port

	log.Printf("Starting Cockaigne server (on %s) ...", hostAndPort)

	persistence.ConnectToDb()
	persistence.DB.AutoMigrate(&auth.Account{}, &games.GameMetadata{})

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./static")

	auth.Register(app)
	header.Register(app)
	home.Register(app)
	account.Register(app)
	games.Register(app)
	highscores.Register(app)

	log.Fatal(app.Listen(hostAndPort))
}
