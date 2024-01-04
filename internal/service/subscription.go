package service

import (
	"github.com/bitstorm-tech/cockaigne/internal/model"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/google/uuid"
)

func CreateSubscription(accountId string, planId int) error {
	activationCode, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	_, err = persistence.DB.Exec(
		"insert into subscriptions (account_id, plan_id, state, activation_code) values ($1, $2, $3, $4)",
		accountId,
		planId,
		model.SubWaitingForActivation,
		activationCode,
	)

	return err
}

func GetPlan(id int) (model.Plan, error) {
	var plan model.Plan
	err := persistence.DB.Get(&plan, "select * from plans where id = $1", id)
	if err != nil {
		return model.Plan{}, err
	}

	return plan, nil
}
