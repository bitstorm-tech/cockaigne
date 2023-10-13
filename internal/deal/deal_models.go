package deal

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Deal struct {
	gorm.Model
	ID              uuid.UUID `gorm:"not null; default: gen_random_uuid(); type: uuid"`
	DealerId        uuid.UUID `gorm:"not null; default: null; type: uuid"`
	Title           string    `gorm:"not null; default: null"`
	Description     string    `gorm:"not null; default: null"`
	CategoryId      int       `gorm:"not null; default: null"`
	DurationInHours int       `gorm:"not null; default: null"`
	Start           time.Time `gorm:"not null; default: null; type: timestamp with time zone"`
	IsTemplate      bool      `gorm:"not null; default: false"`
}

func NewDeal() Deal {
	return Deal{
		Title:           "",
		Description:     "",
		CategoryId:      0,
		DurationInHours: 0,
		Start:           time.Now().Add(1 * time.Hour),
		IsTemplate:      false,
	}
}

func NewDealFromRequest(c *fiber.Ctx) (Deal, string) {
	title := c.FormValue("title")
	if len(title) == 0 {
		return Deal{}, "Bitte einen Titel angeben"
	}

	description := c.FormValue("description")
	if len(description) == 0 {
		return Deal{}, "Bitte eine Beschreibung angeben"
	}

	categoryId, err := strconv.Atoi(c.FormValue("category"))
	if err != nil {
		return Deal{}, "Bitte eine Kategorie auswählen"
	}

	startDate := time.Now()

	if c.FormValue("startInstantly") == "" {
		startDate, err = time.Parse("2006-01-02T15:04", c.FormValue("startDate"))
		if err != nil {
			return Deal{}, "Bitte ein (gültiges) Startdatum angeben"
		}
	}

	duration := 0
	endDate, err := time.Parse("2006-01-02", c.FormValue("endDate"))
	if err == nil {
		duration = int(endDate.Sub(startDate.Truncate(24 * time.Hour)).Hours())
	} else {
		duration, err = strconv.Atoi(c.FormValue("duration"))
		if err != nil {
			return Deal{}, "Bitte entweder eine Laufzeit oder ein Enddatum angeben"
		}
		duration *= 24
	}

	if duration <= 0 {
		return Deal{}, "Das Startdatum muss vor dem Enddatum liegen"
	}

	return Deal{
		Title:           title,
		Description:     description,
		CategoryId:      categoryId,
		Start:           startDate,
		DurationInHours: duration,
	}, ""
}

type Category struct {
	gorm.Model
	Name   string `gorm:"not null; default: null"`
	Active bool   `gorm:"not null; default: true"`
}

type ActiveDeals struct {
	ID                uuid.UUID `gorm:"type: uuid"`
	DealerId          uuid.UUID `gorm:"type: uuid"`
	Title             string
	Description       string
	CategoryId        int
	DurationInMinutes int
	Start             time.Time
	Username          string
	Location          string
	Likes             int
}
