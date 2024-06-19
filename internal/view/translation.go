package view

import (
	"github.com/bitstorm-tech/cockaigne/internal/service"
)

func t(key string, lang string) string {
	return service.T(key, lang)
}
