package view

import (
	"github.com/bitstorm-tech/cockaigne/internal/service"
	"go.uber.org/zap"
)

func t(key string, lang string) string {
	text := service.T(key, lang)
	if len(text) == 0 {
		zap.L().Sugar().Errorf("can't find translation for key=%s / lang=%s", key, lang)
		return "#MISSING_TRANSLATION#"
	}

	return text
}
