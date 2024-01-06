package handler

import (
	"encoding/json"
	"io"
	"os"

	"github.com/bitstorm-tech/cockaigne/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
	"go.uber.org/zap"
)

var webhookSecret = os.Getenv("WEBHOOK_SECRET")

func RegisterStripeHandler(e *echo.Echo) {
	e.POST("/api/stripe/payment-succeeded", paymentSucceeded)
}

func paymentSucceeded(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		zap.L().Sugar().Error("can't read body: ", err)
		return err
	}

	event, err := webhook.ConstructEvent(body, c.Request().Header.Get("Stripe-Signature"), webhookSecret)
	if err != nil {
		zap.L().Sugar().Error("can't construct stribe webhook event: ", err)
		return err
	}

	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			zap.L().Sugar().Error("can't unmarschal checkout.session.completed data: ", err)
			return err
		}

		err = service.ActivateSubscription(session.ID, session.Subscription.ID)
		if err != nil {
			zap.L().Sugar().Error("can't activate subscription: ", err)
			return err
		}
	}

	return nil
}
