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
	e.POST("/api/stripe/webhook", processWebhookEvent)
}

func processWebhookEvent(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		zap.L().Sugar().Error("can't read body: ", err)
		return err
	}

	event, err := webhook.ConstructEvent(body, c.Request().Header.Get("Stripe-Signature"), webhookSecret)
	if err != nil {
		zap.L().Sugar().Error("can't construct stripe webhook event: ", err)
		return err
	}

	if event.Type == stripe.EventTypeCustomerSubscriptionCreated {
		return customerSubscriptionCreated(event)
	}

	if event.Type == stripe.EventTypeCustomerSubscriptionDeleted {
		return customerSubscriptionDeleted(event)
	}

	return nil
}

func customerSubscriptionCreated(event stripe.Event) error {
	var session stripe.CheckoutSession
	err := json.Unmarshal(event.Data.Raw, &session)
	if err != nil {
		zap.L().Sugar().Error("can't unmarshal customer.subscription.created data: ", err)
		return err
	}

	err = service.ActivateSubscription(session.ID, session.Subscription.ID)
	if err != nil {
		zap.L().Sugar().Error("can't activate subscription: ", err)
		return err
	}

	return nil
}

func customerSubscriptionDeleted(event stripe.Event) error {
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		zap.L().Sugar().Error("can't unmarshal customer.subscription.delete data: ", err)
		return err
	}

	err = service.CancelSubscription(subscription.ID)
	if err != nil {
		zap.L().Sugar().Errorf("can't cancel subscription (%s): %v", subscription.ID, err)
		return err
	}

	return nil
}
