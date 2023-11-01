package deal

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Deal struct {
	ID              uuid.UUID
	DealerId        uuid.UUID `db:"dealer_id"`
	Title           string
	Description     string
	CategoryId      int `db:"category_id"`
	DurationInHours int `db:"duration_in_hours"`
	Start           time.Time
	IsTemplate      bool `db:"template"`
	Created         time.Time
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

func DealFromRequest(c *fiber.Ctx) (Deal, string) {
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

	isTemplate := c.FormValue("template") == "on"

	return Deal{
		Title:           title,
		Description:     description,
		CategoryId:      categoryId,
		Start:           startDate,
		DurationInHours: duration,
		IsTemplate:      isTemplate,
	}, ""
}

type Category struct {
	ID       int
	Name     string
	IsActive bool `db:"is_active"`
}

func (c Category) IsFavorite(favCategoryIds []int) bool {
	for _, favCategoryId := range favCategoryIds {
		if favCategoryId == c.ID {
			return true
		}
	}

	return false
}

type DealView struct {
	ID              uuid.UUID
	DealerId        uuid.UUID `db:"dealer_id"`
	Title           string
	Description     string
	CategoryId      int `db:"category_id"`
	DurationInHours int `db:"duration_in_hours"`
	Start           time.Time
	StartTime       time.Time `db:"start_time"`
	Username        string
	Location        string
	Likes           int
}
