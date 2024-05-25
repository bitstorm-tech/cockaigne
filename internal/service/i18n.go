package service

import (
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"go.uber.org/zap"
)

type I18n struct {
	Key string
	De  string
	En  string
}

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
		translations[i.Key][LanguageDe] = i.De
		translations[i.Key][LanguageEn] = i.En
	}
}

func T(key string, lang string) string {
	return translations[key][lang]
}
