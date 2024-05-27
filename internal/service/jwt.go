package service

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type User struct {
	ID              uuid.UUID
	IsDealer        bool
	IsBasicUser     bool
	IsProUser       bool
	IsAuthenticated bool
	Language        string
}

const (
	jwtTokenKeySub         = "sub"
	jwtTokenKeyIsDealer    = "isDealer"
	jwtTokenKeyIsBasicUser = "isBasicUser"
)

func CreateJwtToken(id uuid.UUID, isDealer bool, isBasicUser bool) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		jwtTokenKeySub:         id,
		jwtTokenKeyIsDealer:    isDealer,
		jwtTokenKeyIsBasicUser: isBasicUser,
	})

	signedString, err := token.SignedString(jwtSecret)
	if err != nil {
		zap.L().Sugar().Errorf("Can't signe JWT token: %+v", err)
	}

	return signedString
}

func parseJwtToken(c echo.Context) (jwt.MapClaims, error) {
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

func GetUserFromCookie(c echo.Context) (User, error) {
	claims, err := parseJwtToken(c)
	if err != nil {
		return User{
			IsAuthenticated: false,
		}, fmt.Errorf("can't parse JWT: %v", err)
	}

	id, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return User{
			IsAuthenticated: false,
		}, fmt.Errorf("can't parse userId into UUID: %v", err)
	}

	isDealer := false
	if claims[jwtTokenKeyIsDealer] != nil {
		isDealer = claims[jwtTokenKeyIsDealer].(bool)
	}

	isBasicUser := true
	if claims[jwtTokenKeyIsBasicUser] != nil {
		isBasicUser = claims[jwtTokenKeyIsBasicUser].(bool)
	}

	return User{
		ID:              id,
		IsDealer:        isDealer,
		IsBasicUser:     isBasicUser,
		IsProUser:       !isBasicUser,
		IsAuthenticated: true,
	}, nil
}

func GetUserIdFromCookie(c echo.Context) (uuid.UUID, error) {
	token, err := parseJwtToken(c)
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
	_, err := parseJwtToken(c)
	return err == nil
}

func IsDealer(c echo.Context) bool {
	token, err := parseJwtToken(c)
	if err != nil {
		return false
	}

	return token["isDealer"].(bool)
}

func SetJwtCookie(jwtString string, c echo.Context) {
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    jwtString,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
	}
	c.SetCookie(&cookie)
}
