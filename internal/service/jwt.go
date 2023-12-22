package service

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type User struct {
	ID       uuid.UUID
	IsDealer bool
}

func CreateJwtToken(id uuid.UUID, isDealer bool) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub":      id,
		"isDealer": isDealer,
	})

	signedString, err := token.SignedString(jwtSecret)
	if err != nil {
		zap.L().Sugar().Errorf("Can't signe JWT token: %+v", err)
	}

	return signedString
}

func ParseJwtToken(c echo.Context) (jwt.MapClaims, error) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}

func ParseUser(c echo.Context) (User, error) {
	claims, err := ParseJwtToken(c)
	if err != nil {
		return User{}, fmt.Errorf("can't parse JWT: %v", err)
	}

	id, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return User{}, fmt.Errorf("can't parse userId into UUID: %v", err)
	}

	return User{
		ID:       id,
		IsDealer: claims["isDealer"].(bool),
	}, nil
}

func ParseUserId(c echo.Context) (uuid.UUID, error) {
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

func IsAuthenticated(c echo.Context) bool {
	_, err := ParseJwtToken(c)
	return err == nil
}

func IsDealer(c echo.Context) bool {
	token, err := ParseJwtToken(c)
	if err != nil {
		return false
	}

	return token["isDealer"].(bool)
}
