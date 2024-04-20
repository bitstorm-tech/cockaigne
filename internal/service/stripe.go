package service

import (
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

func CreateStripeCheckoutSession(priceId string, domain string, accountId string) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceId),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(domain + "/subscribe-success"),
		CancelURL:  stripe.String(domain + "/subscribe-cancel/" + accountId),
	}

	s, err := session.New(params)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func DeleteNotActivatedSubscription(accountId string) error {
	_, err := persistence.DB.Exec(
		"delete from subscriptions where account_id = $1 and state = $2",
		accountId,
		model.SubWaitingForActivation,
	)
	return err
}

func CancelSubscription(stripeSubscriptionId string) error {
	_, err := persistence.DB.Exec(
		"update subscriptions set state=$1, canceled=$2 where stripe_subscription_id=$3",
		model.SubCanceled,
		time.Now(),
		stripeSubscriptionId,
	)

	return err
}
