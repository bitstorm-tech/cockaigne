package auth

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Username string    `gorm:"not null; unique"`
	Password string    `grom:"not null"`
	Email    string    `gorm:"not null; unique"`
}
