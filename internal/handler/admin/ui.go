package adminhandler

import (
	"github.com/bitstorm-tech/cockaigne/internal/view"
	adminview "github.com/bitstorm-tech/cockaigne/internal/view/admin"
	"github.com/labstack/echo/v4"
)

func RegisterUiHandler(e *echo.Echo) {
	e.GET("/admin-header", getAdminHeader)
}

func getAdminHeader(c echo.Context) error {
	return view.Render(adminview.Header(), c)
}
