package account

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	ID                   uuid.UUID `gorm:"not null; default: gen_random_uuid(); type: uuid"`
	Username             string    `gorm:"not null; default: null; unique"`
	Password             string    `gorm:"not null; default: null"`
	Email                string    `gorm:"not null; default: null; unique"`
	IsDealer             bool      `gorm:"not null; default: false"`
	Street               string    `gorm:"default: null"`
	HouseNumber          string    `gorm:"default: null"`
	City                 string    `gorm:"default: null"`
	ZipCode              int32     `gorm:"default: null"`
	PhoneNumber          string    `gorm:"default: null"`
	TaxId                string    `gorm:"default: null"`
	DefaultCategory      int       `gorm:"default: null"`
	Location             string    `gorm:"default: null; type: geometry(Point,4326)"`
	SearchRadiusInMeters int       `gorm:"default: 250"`
}
