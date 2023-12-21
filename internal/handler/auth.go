package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bitstorm-tech/cockaigne/internal/model"
	"github.com/bitstorm-tech/cockaigne/internal/service"
	"github.com/bitstorm-tech/cockaigne/internal/view"
	"github.com/labstack/echo/v4"

	"golang.org/x/crypto/bcrypt"
)

func RegisterAuthHandlers(e *echo.Echo) {
	e.GET("/login", func(c echo.Context) error {
		return view.Render(view.Login(), c)
	})

	e.GET("/signup", func(c echo.Context) error {
		return view.Render(view.Register(), c)
	})

	e.GET("/logout", logout)
	e.POST("/api/signup", signup)
	e.POST("/api/login", login)

}

func signup(c echo.Context) error {
	request := model.CreateAccountRequest{}
	err := c.Bind(&request)
	if err != nil {
		c.Logger().Errorf("Error while signup: %v", err)
		return view.RenderAlert("Fehler bei der Registrierung, bitte versuche es später nochmal.", c)
	}

	c.Logger().Debugf("Signup attempt: %+v", request.Email)

	accExists, err := service.AccountExists(request.Email, request.Username)

	if err != nil {
		c.Logger().Errorf("can't signup -> don't know if account already exists: %v", err)
		return view.RenderAlert("Leider gibt es aktuell ein technisches Problem, bitte versuche es später noch einmal!", c)
	}

	if accExists {
		return view.RenderAlert("Benutzername oder E-Mail bereits vergeben", c)
	}

	if request.Password != request.PasswordRepeat {
		return view.RenderAlert("Passwort und Wiederholung stimmen nicht überein", c)
	}

	if request.Username == "" {
		return view.RenderAlert("Bitte einen Benutzernamen angeben", c)
	}

	if request.Email == "" || !strings.Contains(request.Email, "@") {
		return view.RenderAlert("Bitte eine gültige E-Mail angeben", c)
	}

	c.Logger().Debugf("New account: %+v", request.Email)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return view.RenderAlert(err.Error(), c)
	}

	acc := request.ToAccount(string(passwordHash))

	if acc.IsDealer {
		postion, err := service.GetPositionFromAddress(acc.City.String, int(acc.ZipCode.Int32), acc.Street.String, acc.HouseNumber.String)
		if err != nil {
			c.Logger().Errorf("Error while getting position from address: %v", err)
			return view.RenderAlert("Die Adresse konnte nicht gefunden werden", c)
		}

		acc.Location = sql.NullString{
			String: postion.ToWkt(),
		}
	}

	if acc.IsDealer {
		err = service.SaveDealer(acc)
	} else {
		err = service.SaveUser(acc)
	}

	if err != nil {
		return view.RenderAlert(err.Error(), c)
	}

	c.Set("HX-Location", "/login")

	return nil
}

func login(c echo.Context) error {
	request := model.LoginRequest{}
	err := c.Bind(&request)

	if err != nil {
		c.Logger().Errorf("Error while signup %v", err)
		return view.RenderAlert("Login gerade nicht möglich, bitte später nochmal versuchen.", c)
	}

	c.Logger().Debugf("Login attempt: %+v", request.Email)

	acc, err := service.GetAccountByEmail(request.Email)
	if err != nil {
		c.Logger().Errorf("can't get account by email (%s): %v", request.Email, err)
		return view.RenderAlert("Benutzername oder Passwort falsch", c)
	}

	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(request.Password))

	if err != nil {
		return view.RenderAlert("Benutzername oder Passwort falsch", c)
	}

	if acc.IsDealer {
		c.Response().Header().Add("HX-Location", fmt.Sprintf("/dealer/%s", acc.ID.String()))
	} else {
		c.Response().Header().Add("HX-Location", "/user")
	}

	jwtString := service.CreateJwtToken(acc.ID, acc.IsDealer)
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    jwtString,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	}
	c.SetCookie(&cookie)

	return nil
}

func logout(c echo.Context) error {
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		HttpOnly: true,
	}
	c.SetCookie(&cookie)

	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}
