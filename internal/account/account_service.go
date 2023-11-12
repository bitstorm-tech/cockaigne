package account

import (
	"fmt"

	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

func GetSearchRadius(userId uuid.UUID) int {
	var searchRadius = 250
	err := persistence.DB.Get(&searchRadius, "select search_radius_in_meters from accounts where id = $1", userId)

	if err != nil {
		log.Errorf("can't get search radius: %v", err)
	}

	return searchRadius
}

func GetFavoriteCategoryIds(userId uuid.UUID) []int {
	var favoriteCategoryIds []int
	err := persistence.DB.Select(
		&favoriteCategoryIds,
		"select category_id from selected_categories where user_id = $1",
		userId,
	)
	if err != nil {
		log.Errorf("can't get favorite categories: %v", err)
		return []int{}
	}

	return favoriteCategoryIds
}

func GetDefaultCategoryId(userId uuid.UUID) int {
	var defaultCategoryId = -1
	err := persistence.DB.Get(&defaultCategoryId, "select default_category from accounts where id = $1", userId.String())
	if err != nil {
		log.Errorf("can't get default category of dealer: %s", userId.String())
	}

	return defaultCategoryId
}

func GetAccount(userId string) (Account, error) {
	account := Account{}
	err := persistence.DB.Get(&account, "select * from accounts where id = $1", userId)
	if err != nil {
		log.Errorf("can't get account: %v", err)
		return Account{}, err
	}

	return account, nil
}

func GetAccountByEmail(email string) (Account, error) {
	account := Account{}
	err := persistence.DB.Get(&account, "select * from accounts where email = $1", email)
	if err != nil {
		log.Errorf("can't get account: %v", err)
		return Account{}, err
	}

	return account, nil
}

func Exists(email string, username string) (bool, error) {
	var count int
	err := persistence.DB.Get(
		&count,
		"select count(*) from accounts where email ilike $1 or username ilike $2",
		email,
		username,
	)

	if err != nil {
		return true, fmt.Errorf("can't get account count: %v", err)
	}

	return count > 0, nil
}

func SaveDealer(acc Account) error {
	_, err := persistence.DB.Exec(
		"insert into accounts (email, password, street, username, default_category, house_number, city, zip, phone, tax_id, location, is_dealer) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, true)",
		acc.Email,
		acc.Password,
		acc.Street.String,
		acc.Username,
		acc.DefaultCategory.Int32,
		acc.HouseNumber.String,
		acc.City.String,
		acc.ZipCode.Int32,
		acc.Phone.String,
		acc.TaxId.String,
		acc.Location.String,
	)

	return err
}

func SaveUser(acc Account) error {
	_, err := persistence.DB.Exec(
		"insert into accounts (email, password, username, age, gender, is_dealer) values ($1, $2, $3, $4, $5, false)",
		acc.Email,
		acc.Password,
		acc.Username,
		acc.Age.Int32,
		acc.Gender.String,
	)

	return err
}

func UpdateSearchRadius(accountId uuid.UUID, radius int) error {
	_, err := persistence.DB.Exec("update accounts set search_radius_in_meters = $1 where id = $2", radius, accountId)
	return err
}

func UpdateSelectedCategories(userId uuid.UUID, categoryIds []int) error {
	_, err := persistence.DB.Exec("delete from selected_categories where user_id = $1", userId)
	if err != nil {
		return fmt.Errorf("can't delete selected categories: %w", err)
	}

	for _, categoryId := range categoryIds {
		_, err = persistence.DB.Exec("insert into selected_categories (user_id, category_id) values ($1, $2)", userId, categoryId)
		if err != nil {
			return fmt.Errorf("can't insert selected categories: %w", err)
		}
	}

	return nil
}

func UpdateUseLocationService(userId string, useLocationService bool) error {
	_, err := persistence.DB.Exec("update accounts set use_location_service = $1 where id = $2", useLocationService, userId)
	if err != nil {
		return fmt.Errorf("can't save use_location_service for user (%s): %v", userId, err)
	}
	return nil
}
