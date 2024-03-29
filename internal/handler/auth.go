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
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"golang.org/x/crypto/bcrypt"
)

func RegisterAuthHandlers(e *echo.Echo) {
	e.GET("/login", func(c echo.Context) error {
		return view.Render(view.Login(), c)
	})

	e.GET("/signup", func(c echo.Context) error {
		return view.Render(view.Signup(), c)
	})

	e.GET("/signup-complete", completeSignup)
	e.GET("/logout", logout)
	e.POST("/api/signup", signup)
	e.POST("/api/login", login)
	e.POST("/api/basic-login", loginAsBasicUser)
}

func completeSignup(c echo.Context) error {
	email := c.QueryParam("email")
	code := c.QueryParam("code")

	return view.Render(view.SignupComplete(email, code), c)
}

func signup(c echo.Context) error {
	request := model.CreateAccountRequest{}
	err := c.Bind(&request)
	if err != nil {
		zap.L().Sugar().Error("can't bind to signup request: ", err)
		return view.RenderAlert("Fehler bei der Registrierung, bitte versuche es später nochmal.", c)
	}

	zap.L().Sugar().Debug("Signup attempt: ", request.Email)

	accExists, err := service.AccountExists(request.Email, request.Username)

	if err != nil {
		zap.L().Sugar().Error("can't signup -> don't know if account already exists: ", err)
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

	if !request.Agb {
		return view.RenderAlert("Bitte AGB und Datenschutzbedingungen lesen und akzeptieren", c)
	}

	zap.L().Sugar().Debug("new account: ", request.Email)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return view.RenderAlert(err.Error(), c)
	}

	acc := request.ToAccount(string(passwordHash))

	if acc.IsDealer {
		postion, err := service.GetPositionFromAddress(acc.City.String, int(acc.ZipCode.Int32), acc.Street.String, acc.HouseNumber.String)
		if err != nil {
			zap.L().Sugar().Error("Error while getting position from address: ", err)
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
		zap.L().Sugar().Error("can't bind to login request: ", err)
		return view.RenderAlert("Login gerade nicht möglich, bitte später nochmal versuchen.", c)
	}

	zap.L().Sugar().Debug("login attempt: ", request.Email)

	acc, err := service.GetAccountByEmail(request.Email)
	if err != nil {
		zap.L().Sugar().Errorf("can't get account by email (%s): %v", request.Email, err)
		return view.RenderAlert("Benutzername oder Passwort falsch", c)
	}

	if !acc.Active {
		zap.L().Sugar().Infof("login attempt from '%s', but not yet activated", acc.Email)
		return view.RenderAlert("Account noch nicht aktiviert!", c)
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

	jwtString := service.CreateJwtToken(acc.ID, acc.IsDealer, false)
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

func loginAsBasicUser(c echo.Context) error {
	userId := uuid.New()
	jwtString := service.CreateJwtToken(userId, false, true)
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    jwtString,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
	}
	c.SetCookie(&cookie)

	c.Response().Header().Add("HX-Location", "/user")

	service.NewBasicUser(userId.String())

	return nil
}

func logout(c echo.Context) error {
	user, _ := service.ParseUser(c)
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		HttpOnly: true,
	}
	c.SetCookie(&cookie)

	if user.IsBasicUser {
		service.DeleteBasicUser(user.ID.String())
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}
