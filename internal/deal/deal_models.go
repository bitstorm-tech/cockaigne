package deal

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Deal struct {
	gorm.Model
	ID             uuid.UUID `gorm:"not null; default: gen_random_uuid(); type: uuid"`
	DealerId       uuid.UUID `gorm:"not null; default: null; type: uuid"`
	Title          string    `gorm:"not null; default: null"`
	Description    string    `gorm:"not null; default: null"`
	CategoryId     int32     `gorm:"not null; default: null"`
	DurationInDays int32     `gorm:"not null; default: null"`
	Start          time.Time `gorm:"not null; default: null; type: timestamp with time zone"`
	IsTemplate     bool      `gorm:"not null; default: false"`
}

func NewDeal() Deal {
	return Deal{
		Title:          "",
		Description:    "",
		CategoryId:     0,
		DurationInDays: 0,
		Start:          time.Now().Add(1 * time.Hour),
		IsTemplate:     false,
	}
}

type Category struct {
	gorm.Model
	Name   string `gorm:"not null; default: null"`
	Active bool   `gorm:"not null; default: true"`
}
