package account

import (
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func GetSearchRadius(userId uuid.UUID) int {
	account := &Account{}
	err := persistence.DB.Find(account, userId).Error

	searchRadius := 250
	if err == nil {
		searchRadius = account.SearchRadiusInMeters
	} else {
		log.Errorf("can't get search radius: %v", err)
	}

	return searchRadius
}

func GetFavoriteCategoryIds(userId uuid.UUID) []int {
	favoriteCategoryIds := []int{}
	err := persistence.DB.Model(&FavoriteCategory{}).Select("category_id").Where("account_id = ?", userId).Find(&favoriteCategoryIds).Error

	if err != nil {
		log.Errorf("can't get favorite categories: %v", err)
		return []int{}
	}

	return favoriteCategoryIds
}

func GetAccount(userId uuid.UUID) (Account, error) {
	account := Account{}
	err := persistence.DB.Find(&account, userId).Error

	if err != nil {
		log.Errorf("can't get account: %v", err)
		return Account{}, err
	}

	return account, nil
}

func UpdateAccount(account Account) error {
	return persistence.DB.Save(&account).Error
}
