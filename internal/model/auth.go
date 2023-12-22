package model

import (
	"database/sql"
	"strconv"

	"go.uber.org/zap"
)

type CreateAccountRequest struct {
	Username       string
	Password       string
	PasswordRepeat string
	Email          string
	IsDealer       bool
	Street         string
	HouseNumber    string
	City           string
	ZipCode        string
	Phone          string
	TaxId          string
	Category       int32
}

func (c CreateAccountRequest) ToAccount(passwordHash string) Account {
	zipCode := 0

	if c.IsDealer {
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
		IsDealer:        c.IsDealer,
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
