package service

import (
	"github.com/bitstorm-tech/cockaigne/internal/model"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
)

func CreateSubscription(accountId string, planId int, stripeSubscriptionId string) error {
	_, err := persistence.DB.Exec(
		"insert into subscriptions (account_id, plan_id, state, stripe_checkout_session_id) values ($1, $2, $3, $4)",
		accountId,
		planId,
		model.SubWaitingForActivation,
		stripeSubscriptionId,
	)

	return err
}

func ActivateSubscription(stripeCheckoutSessionId string, stripeSubscriptionId string) error {
	_, err := persistence.DB.Exec(
		"update subscriptions set stripe_subscription_id = $1, state = $2, activated = now() where stripe_checkout_session_id = $3",
		stripeSubscriptionId,
		model.SubActive,
		stripeCheckoutSessionId,
	)

	return err
}

func GetPlan(id int) (model.Plan, error) {
	var plan model.Plan
	err := persistence.DB.Get(&plan, "select * from plans where id = $1", id)

	return plan, err
}
