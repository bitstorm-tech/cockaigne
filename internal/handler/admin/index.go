package adminhandler

import (
	"github.com/bitstorm-tech/cockaigne/internal/view"
	adminview "github.com/bitstorm-tech/cockaigne/internal/view/admin"
	"github.com/labstack/echo/v4"
)

func RegisterIndexHandler(e *echo.Echo) {
	e.GET("/", getIndexPage)
}

func getIndexPage(c echo.Context) error {
	// _, err := service.ParseUserId(c)
	// if err != nil {
	// 	redirect.Login(c)
	// }

	return view.Render(adminview.Index(), c)
}
