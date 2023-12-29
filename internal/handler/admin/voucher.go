package adminhandler

import (
	adminmodel "github.com/bitstorm-tech/cockaigne/internal/model/admin"
	"github.com/bitstorm-tech/cockaigne/internal/redirect"
	adminservice "github.com/bitstorm-tech/cockaigne/internal/service/admin"
	"github.com/bitstorm-tech/cockaigne/internal/view"
	adminview "github.com/bitstorm-tech/cockaigne/internal/view/admin"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RegisterVoucherHandler(e *echo.Echo) {
	e.GET("/admin-vouchers", getVouchersPage)
	e.GET("/admin-voucher-create", getVoucherCreatePage)
	e.POST("/admin-voucher", saveVoucher)
}

func saveVoucher(c echo.Context) error {
	// if !auth.IsAuthenticated(c) {
	// 	return c.Redirect("/login")
	// }

	voucherRequest := adminmodel.CreateVoucherRequest{}
	err := c.Bind(&voucherRequest)
	if err != nil {
		zap.L().Sugar().Error("can't parse voucher from request body: ", err)
		return view.RenderAlert(err.Error(), c)
	}

	if voucherRequest.Code == "" {
		return view.RenderAlert("Bitte einen Code angeben", c)
	}

	if voucherRequest.Comment == "" {
		return view.RenderAlert("Bitte einen Kommentar angeben", c)
	}

	voucher, err := voucherRequest.ToVoucher()
	if err != nil {
		zap.L().Sugar().Error("can't create voucher from request: ", err)
		return view.RenderAlert(err.Error(), c)
	}

	err = adminservice.CreateVoucher(voucher)

	if err != nil {
		zap.L().Sugar().Error("can't create voucher in DB: ", err)
		return view.RenderAlert(err.Error(), c)
	}

	return redirect.To("/admin-vouchers", c)
}

func getVoucherCreatePage(c echo.Context) error {
	return view.Render(adminview.VoucherCreate(), c)
}

func getVouchersPage(c echo.Context) error {
	vouchers, err := adminservice.GetVouchers()
	if err != nil {
		zap.L().Sugar().Error("can't get vouchers: ", err)
	}
	return view.Render(adminview.Vouchers(vouchers, err != nil), c)
}
