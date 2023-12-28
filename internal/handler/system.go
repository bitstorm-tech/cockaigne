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
	e.GET("/basic-vs-pro", getBasicVsProPage)
	e.GET("/pricing", getPricingPage)
	e.GET("/active-vouchers-card", getActiveVouchers)
	e.POST("/contact", saveContactMessage)
}

func getActiveVouchers(c echo.Context) error {
	return view.Render(view.VoucherCard([]string{"a", "b", "c"}), c)
}

func getPricingPage(c echo.Context) error {
	return view.Render(view.Pricing(false, []string{}), c)
}

func getBasicVsProPage(c echo.Context) error {
	return view.Render(view.BasicVsPro(), c)
}

func saveContactMessage(c echo.Context) error {
	userId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}

	lastMessageYoungerThen5Minutes, err := service.IsLastContactMessageYoungerThen5Minutes(userId.String())
	if err != nil {
		zap.L().Sugar().Error("can't check if last message is younger then 5 minutes: ", err)
	}

	if lastMessageYoungerThen5Minutes {
		return view.RenderAlert("Du kannst uns nur alle 5 Minuten eine neue Nachricht schreiben, bitte versuche es später noch einmal.", c)
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
