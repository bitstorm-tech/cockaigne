package account

import (
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/google/uuid"
)

func GetSearchRadius(userId uuid.UUID) int {
	account := &Account{}
	err := persistence.DB.Find(account, userId).Error

	searchRadius := 250
	if err == nil {
		searchRadius = account.SearchRadiusInMeters
	}

	return searchRadius
}
