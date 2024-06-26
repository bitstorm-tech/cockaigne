package model

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionState string

const (
	SubscriptionStateWaitingForActivation SubscriptionState = "WAITING_FOR_ACTIVATION"
	SubscriptionStateActive               SubscriptionState = "ACTIVE"
	SubscriptionStateCanceled             SubscriptionState = "CANCELED"
)

type Plan struct {
	ID               int
	Name             string
	StripeProductId  string `db:"stripe_product_id"`
	StripePriceId    string `db:"stripe_price_id"`
	FreeDaysPerMonth int    `db:"free_days_per_month"`
	Active           bool
	Created          time.Time
}

type Subscription struct {
	ID                      int
	AccountId               uuid.UUID `db:"account_id"`
	PlanId                  int       `db:"plan_id"`
	StripeSubscriptionId    string    `db:"stripe_subscription_id"`
	StripeCheckoutSessionId string    `db:"stripe_checkout_session_id"`
	State                   SubscriptionState
	Created                 time.Time
	Activated               time.Time
	Canceled                time.Time
}
