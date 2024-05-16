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

var webhookSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")

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

	switch event.Type {
	case stripe.EventTypeCustomerSubscriptionCreated:
		return customerSubscriptionCreated(event)

	case stripe.EventTypeCustomerSubscriptionUpdated:
		return customerSubscriptionUpdated(event)

	case stripe.EventTypeCustomerSubscriptionDeleted:
		return customerSubscriptionDeleted(event)

	case stripe.EventTypeCheckoutSessionCompleted:
		return checkoutSessionCompleted(event)
	}

	return nil
}

func customerSubscriptionCreated(event stripe.Event) error {
	return nil
}

func customerSubscriptionUpdated(event stripe.Event) error {
	var subscription stripe.Subscription
	err := json.Unmarshal(event.Data.Raw, &subscription)
	if err != nil {
		zap.L().Sugar().Error("can't unmarshal customer.subscription.updated data: ", err)
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

func checkoutSessionCompleted(event stripe.Event) error {
	var checkoutSession stripe.CheckoutSession
	err := json.Unmarshal(event.Data.Raw, &checkoutSession)
	if err != nil {
		zap.L().Sugar().Error("can't unmarshal checkout.session.completed data: ", err)
		return err
	}

	if checkoutSession.Subscription == nil {
		return nil
	}

	trackingId := checkoutSession.Metadata[service.StripeMetadataTrackingId]
	subscriptionId := checkoutSession.Subscription.ID
	err = service.ActivateSubscription(trackingId, subscriptionId)
	if err != nil {
		zap.L().Sugar().Errorf("can't update subscription (trackingId=%s, subscriptionId=%s): %v", trackingId, subscriptionId, err)
		return err
	}

	return nil
}
