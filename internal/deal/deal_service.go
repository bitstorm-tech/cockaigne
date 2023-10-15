package deal

import (
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/gofiber/fiber/v2/log"
)

func GetCategories() []Category {
	categories := []Category{}
	err := persistence.DB.Find(&categories).Error
	if err != nil {
		log.Errorf("can't get categories: %v", err)
	}

	return categories
}
