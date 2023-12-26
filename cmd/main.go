package main

import (
	"os"

	"github.com/bitstorm-tech/cockaigne/internal/handler"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewDevelopment())
	zap.ReplaceGlobals(logger)

	persistence.ConnectToDb()
	persistence.InitS3()

	e := echo.New()
	e.Static("/static", "static")

	handler.RegisterAccountHandlers(e)
	handler.RegisterAuthHandlers(e)
	handler.RegisterDealerHandlers(e)
	handler.RegisterDealHandlers(e)
	handler.RegisterIndexHandlers(e)
	handler.RegisterMapHandlers(e)
	handler.RegisterSystemHandler(e)
	handler.RegisterUiHandlers(e)
	handler.RegisterUserHandlers(e)

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	hostAndPort := host + ":" + port

	logger.Fatal(e.Start(hostAndPort).Error())

}
