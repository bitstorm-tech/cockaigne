package service

import (
	"fmt"

	"github.com/bitstorm-tech/cockaigne/internal/persistence"
)

func SaveContactMessage(accountId string, message string) error {
	if len(message) > 1000 {
		return fmt.Errorf("contact message greater then 1000 characters: %d", len(message))
	}

	_, err := persistence.DB.Exec("insert into contact_messages (account_id, message) values ($1, $2)", accountId, message)

	return err
}
