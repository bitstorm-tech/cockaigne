package service

import (
	"os"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

func init() {
	stripe.Key = os.Getenv("STRIPE_PRIVATE_API_KEY")
}

func CreateStripeCheckoutUrl(priceId string, domain string) (string, error) {
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceId),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(domain + "/subscribe-success"),
		CancelURL:  stripe.String(domain + "/pricing"),
	}

	s, err := session.New(params)
	if err != nil {
		return "", err
	}

	return s.URL, nil
}
