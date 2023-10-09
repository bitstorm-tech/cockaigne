package auth

import (
	"fmt"
	"strings"

	"github.com/bitstorm-tech/cockaigne/internal/account"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
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

	var count int64
	persistence.DB.Model(&account.Account{}).Where("email ilike ?", request.Email).Or("username ilike ?", request.Username).Count(&count)
	if count > 0 {
		return c.Render("partials/alert", fiber.Map{"message": "Benutzername oder E-Mail bereits vergeben"})
	}

	if request.Password != request.PasswordRepeat {
		return c.Render("partials/alert", fiber.Map{"message": "Passwort und Wiederholung stimmen nicht überein"})
	}

	if request.Username == "" {
		return c.Render("partials/alert", fiber.Map{"message": "Bitte einen Benutzernamen angeben"})
	}

	if request.Email == "" || !strings.Contains(request.Email, "@") {
		return c.Render("partials/alert", fiber.Map{"message": "Bitte eine gültige E-Mail angeben"})
	}

	log.Debugf("New account: %+v", request.Email)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Render("partials/alert", fiber.Map{"message": err.Error()})
	}

	acc := request.ToAccount(string(passwordHash))

	err = persistence.DB.Create(&acc).Error
	if err != nil {
		return c.Render("partials/alert", fiber.Map{"message": err.Error()})
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

	acc := account.Account{}
	err = persistence.DB.Where("email ilike ?", request.Email).First(&acc).Error

	if err != nil {
		return c.Render("partials/alert", fiber.Map{"message": "Benutzername oder Passwort falsch"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(request.Password))

	if err != nil {
		return c.Render("partials/alert", fiber.Map{"message": "Benutzername oder Passwort falsch"})
	}

	if acc.IsDealer {
		c.Set("HX-Location", fmt.Sprintf("/dealer/%s", acc.ID.String()))
	} else {
		c.Set("HX-Location", "/user")
	}

	jwt := CreateJwtToken(acc)
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    jwt,
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
