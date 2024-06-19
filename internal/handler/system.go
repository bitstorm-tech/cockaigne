package handler

import (
	"errors"
	"strings"

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
	e.GET("/language-set/:lang", setLanguage)
	e.POST("/contact", saveContactMessage)
	e.POST("/voucher-redeem", redeemVoucher)
}

func setLanguage(c echo.Context) error {
	lang := strings.ToLower(c.Param("lang"))

	if lang != service.LanguageCodeDe && lang != service.LanguageCodeEn {
		zap.L().Sugar().Info("can't change language to: ", lang)
		return nil
	}

	service.SetLanguageCookie(lang, c)

	user, err := service.GetUserFromCookie(c)
	if err == nil {
		if user.IsProUser {
			err := service.SaveLanguage(user.ID.String(), lang)
			if err != nil {
				zap.L().Sugar().Errorf("can't update language to %s for account %s: %v", lang, user.ID, err)
			}
		}
	}

	sourceUrl := c.Request().Header.Get("Referer")

	return redirect.To(sourceUrl, c)
}

func redeemVoucher(c echo.Context) error {
	userId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	voucherCode := c.FormValue("voucher-code")
	err = service.RedeemVoucher(userId.String(), voucherCode)
	if err != nil {
		if !errors.Is(err, service.ErrVoucherAlreadyRedeemed) &&
			!errors.Is(err, service.ErrVoucherNotActive) &&
			!errors.Is(err, service.ErrVoucherCannotBeRedeemed) {
			zap.L().Sugar().Errorf("can't redeem voucher '%s': %v", voucherCode, err)
		} else {
			zap.L().Sugar().Infof("can't redeem voucher '%s': %s", voucherCode, err)
		}
		return view.RenderAlertTranslated("alert.can_t_activeate_voucher", c)
	}

	activeVouchers, err := service.GetActiveVouchers(userId.String())
	lang := service.GetLanguageFromCookie(c)

	return view.Render(view.VoucherCard(activeVouchers, err != nil, lang), c)
}

func getActiveVouchers(c echo.Context) error {
	userId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	activeVouchers, err := service.GetActiveVouchers(userId.String())
	if err != nil {
		zap.L().Sugar().Errorf("can't get active vouchers for user '%s': %v", userId, err)
	}

	lang := service.GetLanguageFromCookie(c)

	return view.Render(view.VoucherCard(activeVouchers, err != nil, lang), c)
}

func getPricingPage(c echo.Context) error {
	lang := service.GetLanguageFromCookie(c)
	return view.Render(view.Pricing(false, []string{}, lang), c)
}

func getBasicVsProPage(c echo.Context) error {
	return view.Render(view.BasicVsPro(), c)
}

func saveContactMessage(c echo.Context) error {
	userId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	lastMessageYoungerThen5Minutes, err := service.IsLastContactMessageYoungerThen5Minutes(userId.String())
	if err != nil {
		zap.L().Sugar().Error("can't check if last message is younger then 5 minutes: ", err)
	}

	if lastMessageYoungerThen5Minutes {
		return view.RenderAlertTranslated("alert.message_delay", c)
	}

	message := c.FormValue("message")
	err = service.SaveContactMessage(userId.String(), message)
	if err != nil {
		zap.L().Sugar().Error("can't save contact message: ", err)
		return view.RenderAlertTranslated("alert.can_t_save_message", c)
	}

	return view.RenderToast("Wir haben deine Nachricht erhalten. Vielen Dank dafür!", c)
}

func getContactPage(c echo.Context) error {
	_, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	lang := service.GetLanguageFromCookie(c)

	return view.Render(view.Contact(lang), c)
}
