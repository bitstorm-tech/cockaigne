package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/bitstorm-tech/cockaigne/internal/model"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailAlreadyExists = errors.New("email already exists")

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
		"select *, st_x(location) || ',' || st_y(location) as location from accounts where email ilike $1",
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

func UpdatePhone(userId string, phone string) error {
	_, err := persistence.DB.Exec("update accounts set phone = $1 where id = $2", phone, userId)
	return err
}

func UpdateTaxId(userId string, taxId string) error {
	_, err := persistence.DB.Exec("update accounts set tax_id = $1 where id = $2", taxId, userId)
	return err
}

func UpdateDefaultCategory(userId string, categoryId int) error {
	_, err := persistence.DB.Exec("update accounts set default_category = $1 where id = $2", categoryId, userId)
	return err
}

func UpdateDealerAddress(dealerId string, street string, houseNumber string, city string, zip int32) error {
	_, err := persistence.DB.Exec(
		"update accounts set street = $1, house_number = $2, city = $3, zip = $4 where id = $5",
		street,
		houseNumber,
		city,
		zip,
		dealerId,
	)
	return err
}

func GetProfileImage(accountId string) (string, error) {
	imageUrl, err := persistence.GetImageUrl(persistence.ProfileImagesFolder + "/" + accountId)
	if err != nil {
		return "", err
	}

	return imageUrl, nil
}

func UsernameExists(accountId string, username string) bool {
	exists := true
	err := persistence.DB.Get(
		&exists,
		"select exists (select * from accounts where username ilike $1 and id != $2)",
		username,
		accountId,
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

func SetActivationCode(email string, code int) error {
	_, err := persistence.DB.Exec(
		"update accounts set activation_code = $1 where email ilike $2",
		code,
		email,
	)

	return err
}

func ActivateAccount(code int) error {
	_, err := persistence.DB.Exec("update accounts set active = true where activation_code = $1", code)

	return err
}

func PreparePasswordChange(email string, accountId string, baseUrl string) error {
	code, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	if len(email) > 0 {
		_, err := persistence.DB.Exec(
			"update accounts set change_password_code = $1 where email ilike $2",
			code.String(),
			email,
		)
		if err != nil {
			return err
		}
	} else {
		err := persistence.DB.Get(
			&email,
			"update accounts set change_password_code = $1 where id = $2 returning email",
			code.String(),
			accountId,
		)
		if err != nil {
			return err
		}
	}

	return SendPasswordChangeEmail(email, code.String(), baseUrl)
}

func ChangePassword(code string, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = persistence.DB.Exec(
		"update accounts set password = $1, change_password_code = null where change_password_code = $2",
		passwordHash,
		code,
	)

	return err
}

func PrepareEmailChange(accountId string, newEmail string, baseUrl string) error {
	code, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	var emailAlreadyExists bool
	err = persistence.DB.Get(&emailAlreadyExists, "select exists(select * from accounts where email ilike $1)", newEmail)
	if err != nil {
		return err
	}

	if emailAlreadyExists {
		return ErrEmailAlreadyExists
	}

	_, err = persistence.DB.Exec(
		"update accounts set change_email_code = $1, new_email = $2 where id = $3",
		code.String(),
		newEmail,
		accountId,
	)
	if err != nil {
		return err
	}

	return SendEmailChangeEmail(newEmail, code.String(), baseUrl)
}

func ChangeEmail(code string) error {
	_, err := persistence.DB.Exec(
		"update accounts set email = new_email, new_email = null, change_email_code = null where change_email_code = $1",
		code,
	)

	return err
}
