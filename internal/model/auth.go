package model

import (
	"database/sql"
	"strconv"

	"go.uber.org/zap"
)

type CreateAccountRequest struct {
	Username       string `form:"username"`
	Password       string `form:"password"`
	PasswordRepeat string `form:"passwordRepeat"`
	Email          string `form:"email"`
	IsDealer       string `form:"isDealer"`
	Street         string `form:"street"`
	HouseNumber    string `form:"houseNumber"`
	City           string `form:"city"`
	ZipCode        string `form:"zipCode"`
	Phone          string `form:"phone"`
	TaxId          string `form:"taxId"`
	Category       int32  `form:"category"`
	Age            string `form:"age"`
	Gender         string `form:"gender"`
	Agb            string `form:"agb"`
}

func (c CreateAccountRequest) ToAccount(passwordHash string) Account {
	zipCode := 0
	isDealer := c.IsDealer == "on"

	if isDealer {
		var err error
		zipCode, err = strconv.Atoi(c.ZipCode)

		if err != nil {
			zap.L().Sugar().Errorf("Can't convert zip code: %+v", err)
		}
	}

	return Account{
		Username:        c.Username,
		Password:        passwordHash,
		Email:           c.Email,
		IsDealer:        isDealer,
		Street:          sql.NullString{String: c.Street},
		HouseNumber:     sql.NullString{String: c.HouseNumber},
		City:            sql.NullString{String: c.City},
		ZipCode:         sql.NullInt32{Int32: int32(zipCode)},
		Phone:           sql.NullString{String: c.Phone},
		TaxId:           sql.NullString{String: c.TaxId},
		DefaultCategory: sql.NullInt32{Int32: c.Category},
	}
}

type LoginRequest struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}
