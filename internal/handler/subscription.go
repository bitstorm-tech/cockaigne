package handler

import (
	"fmt"
	"strconv"

	"github.com/bitstorm-tech/cockaigne/internal/redirect"
	"github.com/bitstorm-tech/cockaigne/internal/service"
	"github.com/bitstorm-tech/cockaigne/internal/view"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RegisterSubscriptionHandler(e *echo.Echo) {
	e.POST("/subscribe/:planId", subscribe)
	e.GET("/subscription-activate/:accountId/:activationCode", activateSubscription)
}

func subscribe(c echo.Context) error {
	accountId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}

	planIdString := c.Param("planId")
	planId, err := strconv.Atoi(planIdString)
	if err != nil {
		zap.L().Sugar().Errorf("can't convert plan ID '%s' to integer: %v", planIdString, err)
		return view.RenderAlert("Momentan können keine Abonemants abgeschlossen werden. Bitte versuche es später nochmal.", c)
	}

	err = service.CreateSubscription(accountId.String(), planId)
	if err != nil {
		zap.L().Sugar().Error("can't create subscription: ", err)
		return view.RenderAlert("Momentan können keine Abonemants abgeschlossen werden. Bitte versuche es später nochmal.", c)
	}

	plan, err := service.GetPlan(planId)
	if err != nil {
		zap.L().Sugar().Error("can't get plan: ", err)
		return view.RenderAlert("Momentan können keine Abonemants abgeschlossen werden. Bitte versuche es später nochmal.", c)
	}

	domain := fmt.Sprintf("%s://%s", c.Scheme(), c.Request().Host)
	url, err := service.CreateStripeCheckoutUrl(plan.StripePriceId, domain)
	if err != nil {
		zap.L().Sugar().Error("can't create stripe checkout url: ", err)
		return view.RenderAlert("Momentan können keine Abonemants abgeschlossen werden. Bitte versuche es später nochmal.", c)
	}

	c.Response().Header().Add("HX-Redirect", url)

	return nil
}

func activateSubscription(c echo.Context) error {
	return nil
}
