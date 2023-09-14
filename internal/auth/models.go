package auth

import (
	"strconv"

	"github.com/bitstorm-tech/cockaigne/internal/account"
	"github.com/gofiber/fiber/v2/log"
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
	PhoneNumber    string
	TaxId          string
}

func (c CreateAccountRequest) ToAccount(passwordHash string) account.Account {
	zipCode := 0

	if c.IsDealer {
		var err error
		zipCode, err = strconv.Atoi(c.ZipCode)

		if err != nil {
			log.Errorf("Can't convert zip code: %+v", err)
		}
	}

	return account.Account{
		Username:    c.Username,
		Password:    string(passwordHash),
		Email:       c.Email,
		IsDealer:    c.IsDealer,
		Street:      c.Street,
		HouseNumber: c.HouseNumber,
		City:        c.City,
		ZipCode:     int32(zipCode),
		PhoneNumber: c.PhoneNumber,
		TaxId:       c.TaxId,
	}
}

type LoginRequest struct {
	Email    string
	Password string
}
