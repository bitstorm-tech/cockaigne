package auth

import (
	"os"

	"github.com/bitstorm-tech/cockaigne/internal/account"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func CreateJwtToken(acc account.Account) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub":      acc.ID,
		"isDealer": acc.IsDealer,
	})

	signedString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Errorf("Can't signe JWT token: %+v", err)
	}

	return signedString
}

func ParseJwtToken(c *fiber.Ctx) (jwt.MapClaims, error) {
	tokenString := c.Cookies("jwt")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}

func ParseUserId(c *fiber.Ctx) (uuid.UUID, error) {
	token, err := ParseJwtToken(c)
	if err != nil {
		return uuid.Nil, err
	}

	userId, err := uuid.Parse(token["sub"].(string))
	if err != nil {
		return uuid.Nil, err
	}

	return userId, nil
}

func IsAuthenticated(c *fiber.Ctx) bool {
	_, err := ParseJwtToken(c)
	return err == nil
}
