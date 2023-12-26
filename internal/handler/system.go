package handler

import (
	"github.com/bitstorm-tech/cockaigne/internal/redirect"
	"github.com/bitstorm-tech/cockaigne/internal/service"
	"github.com/bitstorm-tech/cockaigne/internal/view"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RegisterSystemHandler(e *echo.Echo) {
	e.GET("/contact", getContactPage)
	e.POST("/contact", saveContactMessage)
}

func saveContactMessage(c echo.Context) error {
	userId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}

	message := c.FormValue("message")
	err = service.SaveContactMessage(userId.String(), message)
	if err != nil {
		zap.L().Sugar().Error("can't save contact message: ", err)
		return view.RenderAlert("Leider ist beim speichern deiner Nachricht etwas schief gegangen, bitte versuche es später noch einmal.", c)
	}

	return view.RenderToast("Wir haben deine Nachricht erhalten. Vielen Dank dafür!", c)
}

func getContactPage(c echo.Context) error {
	return view.Render(view.Contact(), c)
}
