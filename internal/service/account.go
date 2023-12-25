package service

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/bitstorm-tech/cockaigne/internal/model"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func GetSearchRadius(userId uuid.UUID) int {
	var searchRadius = 250
	err := persistence.DB.Get(&searchRadius, "select search_radius_in_meters from accounts where id = $1", userId)

	if err != nil {
		zap.L().Sugar().Errorf("can't get search radius: %v", err)
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
		zap.L().Sugar().Errorf("can't get favorite categories: %v", err)
		return []int{}
	}

	return favoriteCategoryIds
}

func GetDefaultCategoryId(userId uuid.UUID) int {
	var defaultCategoryId = -1
	err := persistence.DB.Get(&defaultCategoryId, "select default_category from accounts where id = $1", userId.String())
	if err != nil {
		zap.L().Sugar().Errorf("can't get default category of dealer: %s", userId.String())
	}

	return defaultCategoryId
}

func GetAccount(userId string) (model.Account, error) {
	account := model.Account{}
	err := persistence.DB.Get(
		&account,
		"select *, st_x(location) || ',' || st_y(location) as location from accounts where id = $1",
		userId,
	)
	if err != nil {
		zap.L().Sugar().Errorf("can't get account: %v", err)
		return model.Account{}, err
	}

	return account, nil
}

func GetAccountByEmail(email string) (model.Account, error) {
	var account model.Account
	err := persistence.DB.Get(
		&account,
		"select *, st_x(location) || ',' || st_y(location) as location from accounts where email = $1",
		email,
	)
	if err != nil {
		zap.L().Sugar().Errorf("can't get account: %v", err)
		return model.Account{}, err
	}

	return account, nil
}

func AccountExists(email string, username string) (bool, error) {
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

func SaveDealer(acc model.Account) error {
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

func SaveUser(acc model.Account) error {
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
	return err
}

func UpdateLocation(userId string, location model.Point) error {
	_, err := persistence.DB.Exec("update accounts set location = $1 where id = $2", location.ToWkt(), userId)
	return err
}

func UpdateUsername(userId string, username string) error {
	_, err := persistence.DB.Exec("update accounts set username = $1 where id = $2", username, userId)
	return err
}

func GetProfileImage(accountId string) (string, error) {
	imageUrl, err := persistence.GetImageUrl(persistence.ProfileImagesFolder + "/" + accountId)
	if err != nil {
		return "", err
	}

	return imageUrl, nil
}

func UsernameExists(username string) bool {
	exists := true
	err := persistence.DB.Get(
		&exists,
		"select exists (select * from accounts where username ilike $1)",
		username,
	)
	if err != nil {
		zap.L().Sugar().Errorf("can't check if username '%s' already exists: %v", username, err)
	}

	return exists
}

func SaveProfileImage(accountId string, image *multipart.FileHeader) (string, error) {
	tokens := strings.Split(image.Filename, ".")
	fileExtension := tokens[len(tokens)-1]
	path := fmt.Sprintf("%s/%s.%s", persistence.ProfileImagesFolder, accountId, fileExtension)
	err := persistence.UploadImage(path, image)
	if err != nil {
		return "", err
	}

	imageUrl := persistence.S3BaseUrl + "/" + path

	return imageUrl, nil
}

func DeleteProfileImage(accountId string) error {
	path, err := persistence.GetImageUrl(fmt.Sprintf("%s/%s", persistence.ProfileImagesFolder, accountId))
	if err != nil {
		return err
	}

	path = strings.ReplaceAll(path, persistence.S3BaseUrl+"/", "")

	return persistence.DeleteImage(path)
}
