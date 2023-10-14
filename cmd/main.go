package main

import (
	"os"
	"strings"

	"github.com/bitstorm-tech/cockaigne/internal/account"
	"github.com/bitstorm-tech/cockaigne/internal/auth"
	"github.com/bitstorm-tech/cockaigne/internal/deal"
	"github.com/bitstorm-tech/cockaigne/internal/dealer"
	"github.com/bitstorm-tech/cockaigne/internal/home"
	"github.com/bitstorm-tech/cockaigne/internal/like"
	"github.com/bitstorm-tech/cockaigne/internal/maps"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/bitstorm-tech/cockaigne/internal/ui"
	"github.com/bitstorm-tech/cockaigne/internal/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/template/html/v2"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	hostAndPort := host + ":" + port

	log.Debugf("Starting Cockaigne server (on %s) ...", hostAndPort)

	persistence.ConnectToDb()

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./static")

	auth.Register(app)
	ui.Register(app)
	home.Register(app)
	account.Register(app)
	user.Register(app)
	dealer.Register(app)
	deal.Register(app)
	maps.Register(app)

	log.Fatal(app.Listen(hostAndPort))
}
