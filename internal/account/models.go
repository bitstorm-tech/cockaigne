package account

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	ID       uuid.UUID `gorm:"not null; default: gen_random_uuid(); type: uuid"`
	Username string    `gorm:"not null; default: null; unique"`
	Password string    `gorm:"not null; default: null"`
	Email    string    `gorm:"not null; default: null; unique"`
	IsDealer bool      `gorm:"not null; default: false"`
}
