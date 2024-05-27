package service

import (
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type I18n struct {
	Key string
	De  string
	En  string
}

const (
	LanguageCodeDe = "de"
	LanguageCodeEn = "en"
)

const cookieLanguageKey = "lang"

var translations map[string]map[string]string

func LoadI18n() {
	zap.L().Sugar().Info("loading translations from i18n table ...")
	var i18n []I18n
	err := persistence.DB.Select(
		&i18n,
		"select * from i18n",
	)
	if err != nil {
		zap.L().Sugar().Error("can't load translations from database: ", err)
		return
	}

	translations = map[string]map[string]string{}

	for _, i := range i18n {
		translations[i.Key] = map[string]string{}
		translations[i.Key][LanguageCodeDe] = i.De
		translations[i.Key][LanguageCodeEn] = i.En
	}
}

func T(key string, lang string) string {
	return translations[key][lang]
}

func GetLanguageFromCookie(c echo.Context) string {
	cookie, err := c.Request().Cookie(cookieLanguageKey)
	if err != nil {
		if err == http.ErrNoCookie {

		}
		zap.L().Sugar().Error("can't get language cookie: ", err)
		return LanguageCodeDe
	}

	lang := strings.ToLower(cookie.Value)
	if lang == LanguageCodeDe || lang == LanguageCodeEn {
		return lang
	}

	user, err := GetUserFromCookie(c)
	if err != nil {
		zap.L().Sugar().Error("can't get user cookie: ", err)
		return LanguageCodeDe
	}

	if user.IsProUser {
		lang, err := GetLanguage(user.ID.String())
		if err != nil {
			zap.L().Sugar().Errorf("can't get language of user %s, %v", user.ID, err)
			return LanguageCodeDe
		}

		return lang
	}

	return LanguageCodeDe
}

func SetLanguageCookie(lang string) {}
