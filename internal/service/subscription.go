package service

import (
	"fmt"
	"time"

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

func GetFreeDaysLeftFromSubscription(dealerId string) (int, error) {
	daysUsed := 0
	err := persistence.DB.Get(
		&daysUsed,
		`select coalesce(sum(duration_in_hours) / 24, 0)
		from active_deals_view 
		where date_trunc('month', start) = date_trunc('month', now()) 
		  and date_trunc('year', start) = date_trunc('year', now()) 
		  and dealer_id = $1`,
		dealerId,
	)
	if err != nil {
		return 0, err
	}

	freeDaysPerMont, err := GetFreeDaysPerMonthFromSubscription(dealerId)
	if err != nil {
		return 0, err
	}

	return freeDaysPerMont - daysUsed, nil
}

func GetFreeDaysPerMonthFromSubscription(dealerId string) (int, error) {
	freeDaysPerMont := 0
	err := persistence.DB.Get(
		&freeDaysPerMont,
		"select free_days_per_month from subscriptions s join plans p on s.plan_id = p.id where state = $1 and account_id = $2",
		model.SubActive,
		dealerId,
	)
	if err != nil {
		return 0, err
	}

	return freeDaysPerMont, nil
}

func GetSubscriptionPeriodEndDate(dealerId string) (string, error) {
	var activationDate time.Time
	err := persistence.DB.Get(
		&activationDate,
		"select activated from subscriptions where account_id = $1",
		dealerId,
	)
	if err != nil {
		return "", err
	}

	activationDay := activationDate.Day()
	nextMonth := time.Now().Month() + 1
	currentYear := time.Now().Year()

	if nextMonth == 13 {
		nextMonth = 1
		currentYear = currentYear + 1
	}

	return fmt.Sprintf("%02d.%02d.%d", activationDay, nextMonth, currentYear), nil
}
