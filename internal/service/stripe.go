package service

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"time"

	"github.com/bitstorm-tech/cockaigne/internal/model"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

func init() {
	stripe.Key = os.Getenv("STRIPE_PRIVATE_API_KEY")
}

const StripeMetadataTrackingId = "trackingId"

func CreateStripeCheckoutSessionForSubscription(priceId string, domain string) (*stripe.CheckoutSession, error) {
	trackingId := uuid.New().String()
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceId),
				Quantity: stripe.Int64(1),
			},
		},
		Metadata: map[string]string{
			StripeMetadataTrackingId: trackingId,
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(domain + "/subscribe-success"),
		CancelURL:  stripe.String(domain + "/subscribe-cancel/" + trackingId),
	}

	s, err := session.New(params)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func createStripeCheckoutSessionForDynamicPrice(amount int64, days int, domain string, dealId string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String("EUR"),
					UnitAmount: stripe.Int64(amount),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(fmt.Sprintf("Preis für %d Tag(e) Laufzeit", days)),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(domain + "/deal-payed/" + dealId),
		CancelURL:  stripe.String(domain),
	}

	s, err := session.New(params)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func DeleteNotActivatedSubscription(trackingId string) error {
	_, err := persistence.DB.Exec(
		"delete from subscriptions where stripe_tracking_id = $1 and state = $2",
		trackingId,
		model.SubscriptionStateWaitingForActivation,
	)
	return err
}

func CancelSubscription(stripeSubscriptionId string) error {
	_, err := persistence.DB.Exec(
		"update subscriptions set state=$1, canceled=$2 where stripe_subscription_id=$3",
		model.SubscriptionStateCanceled,
		time.Now(),
		stripeSubscriptionId,
	)

	return err
}

func DoStripePayment(dealerId string, dealId string, dealDays int, domain string) (*stripe.CheckoutSession, error) {
	hasActiveSub, err := HasActiveSubscription(dealerId)
	if err != nil {
		return nil, err
	}

	if hasActiveSub {
		freeDaysLeft, err := GetFreeDaysLeftFromSubscription(dealerId)
		if err != nil {
			return nil, err
		}

		if freeDaysLeft >= dealDays {
			return nil, nil
		}

		daysToPay := dealDays - freeDaysLeft

		amount := int64(daysToPay * 499)
		return createStripeCheckoutSessionForDynamicPrice(amount, daysToPay, domain, dealId)
	}

	discount, err := GetHighestVoucherDiscount(dealerId)
	if err != nil {
		return nil, err
	}

	amount := int64(dealDays * 499)
	amount = (amount * int64(100-discount)) / 100

	if amount > 0 {
		return createStripeCheckoutSessionForDynamicPrice(amount, dealDays, domain, dealId)
	}

	return nil, nil
}
