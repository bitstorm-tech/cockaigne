package auth

import (
	"database/sql"
	"fmt"
	"github.com/bitstorm-tech/cockaigne/internal/account"
	"github.com/bitstorm-tech/cockaigne/internal/ui"
	"strings"

	"github.com/bitstorm-tech/cockaigne/internal/auth/jwt"
	"github.com/bitstorm-tech/cockaigne/internal/geo"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"
)

func Register(app *fiber.App) {
	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("pages/login", nil, "layouts/main")
	})

	app.Get("/signup", func(c *fiber.Ctx) error {
		return c.Render("pages/signup", nil, "layouts/main")
	})

	app.Post("/api/signup", signup)

	app.Post("/api/login", login)

	app.Get("/logout", logout)
}

func signup(c *fiber.Ctx) error {
	request := CreateAccountRequest{}
	err := c.BodyParser(&request)

	if err != nil {
		log.Errorf("Error while signup: %v", err)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	log.Debugf("Signup attempt: %+v", request.Email)

	accExists, err := account.Exists(request.Email, request.Username)

	if err != nil {
		log.Errorf("can't signup -> don't know if account already exists: %v", err)
		return ui.ShowAlert(c, "Leider gibt es aktuell ein technisches Problem. Bitte später noch einmal versuchen!")
	}

	if accExists {
		return ui.ShowAlert(c, "Benutzername oder E-Mail bereits vergeben")
	}

	if request.Password != request.PasswordRepeat {
		return ui.ShowAlert(c, "Passwort und Wiederholung stimmen nicht überein")
	}

	if request.Username == "" {
		return ui.ShowAlert(c, "Bitte einen Benutzernamen angeben")
	}

	if request.Email == "" || !strings.Contains(request.Email, "@") {
		return ui.ShowAlert(c, "Bitte eine gültige E-Mail angeben")
	}

	log.Debugf("New account: %+v", request.Email)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return ui.ShowAlert(c, err.Error())
	}

	acc := request.ToAccount(string(passwordHash))

	if acc.IsDealer {
		postion, err := geo.GetPositionFromAddress(acc.City.String, int(acc.ZipCode.Int32), acc.Street.String, acc.HouseNumber.String)
		if err != nil {
			log.Errorf("Error while getting position from address: %v", err)
			return ui.ShowAlert(c, "Die Adresse konnte nicht gefunden werden")
		}

		acc.Location = sql.NullString{
			String: postion.ToWkt(),
		}
	}

	if acc.IsDealer {
		err = account.SaveDealer(acc)
	} else {
		err = account.SaveUser(acc)
	}

	if err != nil {
		return ui.ShowAlert(c, err.Error())
	}

	c.Set("HX-Location", "/login")

	return c.SendStatus(200)
}

func login(c *fiber.Ctx) error {
	request := LoginRequest{}
	err := c.BodyParser(&request)

	if err != nil {
		log.Errorf("Error while signup %v", err)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	log.Debugf("Login attempt: %+v", request.Email)

	acc, err := account.GetAccountByEmail(request.Email)
	if err != nil {
		log.Errorf("can't get account by email (%s): %v", request.Email, err)
		return ui.ShowAlert(c, "Benutzername oder Passwort falsch")
	}

	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(request.Password))

	if err != nil {
		return ui.ShowAlert(c, "Benutzername oder Passwort falsch")
	}

	if acc.IsDealer {
		c.Set("HX-Location", fmt.Sprintf("/dealer/%s", acc.ID.String()))
	} else {
		c.Set("HX-Location", "/user")
	}

	jwtString := jwt.CreateJwtToken(acc.ID, acc.IsDealer)
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    jwtString,
		HTTPOnly: true,
	})

	return nil
}

func logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		HTTPOnly: true,
	})

	return c.Redirect("/login")
}
