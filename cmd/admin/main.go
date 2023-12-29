package main

import (
	"os"

	"github.com/bitstorm-tech/cockaigne/internal/handler"
	adminhandler "github.com/bitstorm-tech/cockaigne/internal/handler/admin"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {
	logger := zap.Must(zap.NewDevelopment())
	zap.ReplaceGlobals(logger)

	persistence.ConnectToDb()

	e := echo.New()
	e.Static("/static", "static")

	adminhandler.RegisterIndexHandler(e)
	adminhandler.RegisterUiHandler(e)
	adminhandler.RegisterVoucherHandler(e)
	handler.RegisterUiHandlers(e)

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	hostAndPort := host + ":" + port

	logger.Fatal(e.Start(hostAndPort).Error())

}
