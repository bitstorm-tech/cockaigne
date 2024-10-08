package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/bitstorm-tech/cockaigne/internal/redirect"

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
		lang := service.GetLanguageFromCookie(c)
		return view.Render(view.Login(lang), c)
	})

	e.GET("/signup", func(c echo.Context) error {
		lang := service.GetLanguageFromCookie(c)
		return view.Render(view.Signup(lang), c)
	})

	e.GET("/logout", logout)
	e.GET("/signup-complete", completeSignup)
	e.POST("/api/basic-login", loginAsBasicUser)
	e.POST("/api/login", login)
	e.POST("/api/signup", signup)
}

func completeSignup(c echo.Context) error {
	email := c.QueryParam("email")
	code := c.QueryParam("code")
	lang := service.GetLanguageFromCookie(c)

	return view.Render(view.SignupComplete(email, code, lang), c)
}

func signup(c echo.Context) error {
	request := model.CreateAccountRequest{}
	err := c.Bind(&request)
	if err != nil {
		zap.L().Sugar().Error("can't bind to signup request: ", err)
		return view.RenderAlertTranslated("alert.error_while_signup", c)
	}

	zap.L().Sugar().Debug("Signup attempt: ", request.Email)

	accExists, err := service.AccountExists(request.Email, request.Username)

	if err != nil {
		zap.L().Sugar().Error("can't signup -> don't know if account already exists: ", err)
		return view.RenderAlertTranslated("alert.technical_problem", c)
	}

	if accExists {
		return view.RenderAlertTranslated("alert.username_or_email_already_used", c)
	}

	if request.Password != request.PasswordRepeat {
		return view.RenderAlertTranslated("alert.password_repeat_not_matching", c)
	}

	if request.Username == "" {
		return view.RenderAlertTranslated("alert.provide_username", c)
	}

	if request.Email == "" || !strings.Contains(request.Email, "@") {
		return view.RenderAlertTranslated("alert.provide_email", c)
	}

	if request.Agb != "on" {
		return view.RenderAlertTranslated("alert.accept_terms_and_privacy", c)
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
			return view.RenderAlertTranslated("alert.invalid_address", c)
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
		zap.L().Sugar().Error("can't save account: ", err)
		return view.RenderAlertTranslated("alert.can_t_create_account", c)
	}

	baseUrl := service.GetBaseUrl(c)
	err = service.SendAccountActivationMail(request.Email, baseUrl)
	if err != nil {
		zap.L().Sugar().Error("can't send account activation email: ", err)
		return view.RenderAlertTranslated("alert.can_t_create_account", c)
	}

	return redirect.To("/signup-complete", c)
}

func login(c echo.Context) error {
	request := model.LoginRequest{}
	err := c.Bind(&request)

	if err != nil {
		zap.L().Sugar().Error("can't bind to login request: ", err)
		return view.RenderAlertTranslated("alert.login_not_possible", c)
	}

	acc, err := service.GetAccountByEmail(request.Email)
	if err != nil {
		zap.L().Sugar().Errorf("can't get account by email (%s): %v", request.Email, err)
		return view.RenderAlertTranslated("alert.invalid_username_or_password", c)
	}

	if !acc.Active {
		zap.L().Sugar().Infof("login attempt from '%s', but not yet activated", acc.Email)
		return view.RenderAlertTranslated("alert.account_not_activated", c)
	}

	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(request.Password))

	if err != nil {
		return view.RenderAlertTranslated("alert.invalid_username_or_password", c)
	}

	if acc.IsDealer {
		c.Response().Header().Add("HX-Location", fmt.Sprintf("/dealer/%s", acc.ID.String()))
	} else {
		c.Response().Header().Add("HX-Location", "/user")
	}

	jwtString := service.CreateJwtToken(acc.ID, acc.IsDealer, false)
	service.SetJwtCookie(jwtString, c)
	return nil
}

func loginAsBasicUser(c echo.Context) error {
	userId := uuid.New()
	jwtString := service.CreateJwtToken(userId, false, true)
	service.SetJwtCookie(jwtString, c)

	c.Response().Header().Add("HX-Location", "/user")

	service.NewBasicUser(userId.String())

	return nil
}

func logout(c echo.Context) error {
	user, _ := service.GetUserFromCookie(c)
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
