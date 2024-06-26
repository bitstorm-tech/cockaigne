package handler

import (
	"net/http"
	"strconv"

	"github.com/bitstorm-tech/cockaigne/internal/redirect"
	"github.com/bitstorm-tech/cockaigne/internal/service"
	"github.com/bitstorm-tech/cockaigne/internal/view"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RegisterSubscriptionHandler(e *echo.Echo) {
	e.POST("/subscribe/:planId", subscribe)
	e.GET("/subscribe-success", subscribeSuccess)
	e.GET("/subscribe-cancel/:trackingId", subscribeCancel)
}

func subscribe(c echo.Context) error {
	accountId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	planIdString := c.Param("planId")
	planId, err := strconv.Atoi(planIdString)
	if err != nil {
		zap.L().Sugar().Errorf("can't convert plan ID '%s' to integer: %v", planIdString, err)
		return view.RenderAlertTranslated("alert.can_t_conclude_subscription", c)
	}

	plan, err := service.GetPlan(planId)
	if err != nil {
		zap.L().Sugar().Errorf("can't get plan (id=%s): %v", planIdString, err)
		return view.RenderAlertTranslated("alert.can_t_conclude_subscription", c)
	}

	baseUrl := service.GetBaseUrl(c)
	checkoutSession, err := service.CreateStripeCheckoutSessionForSubscription(plan.StripePriceId, baseUrl)
	if err != nil {
		zap.L().Sugar().Errorf(
			"can't create stripe checkout session (account=%s, plan=%s, stripePrice=%s, stripeProduct=%s): %v",
			accountId,
			planIdString,
			plan.StripePriceId,
			plan.StripeProductId,
			err,
		)
		return view.RenderAlertTranslated("alert.can_t_conclude_subscription", c)
	}

	err = service.CreateSubscription(accountId.String(), planId, checkoutSession.Metadata[service.StripeMetadataTrackingId])
	if err != nil {
		zap.L().Sugar().Errorf(
			"can't create subscription (account=%s, checkoutSession=%s): %v",
			accountId,
			checkoutSession.ID,
			err,
		)
		return view.RenderAlertTranslated("alert.can_t_conclude_subscription", c)
	}

	c.Response().Header().Add("HX-Redirect", checkoutSession.URL)

	return nil
}

func subscribeSuccess(c echo.Context) error {
	return view.Render(view.SubscribeSuccess(), c)
}

func subscribeCancel(c echo.Context) error {
	trackingId := c.Param("trackingId")
	err := service.DeleteNotActivatedSubscription(trackingId)
	if err != nil {
		zap.L().Sugar().Error("can't delete not activated subscription: ", err)
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/pricing")
}
