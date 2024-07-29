package model

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type DealPaymentState string

const (
	DealPaymentStatePayed   DealPaymentState = "PAYED"
	DealPaymentStatePending DealPaymentState = "PAYMENT_PENDING"
)

type DealState string

const (
	DealStateActive   DealState = "active"
	DealStatePast     DealState = "past"
	DealStateFuture   DealState = "future"
	DealStateTemplate DealState = "template"
)

func ToDealState(state string) DealState {
	switch strings.ToLower(state) {
	case "active":
		return DealStateActive
	case "past":
		return DealStatePast
	case "future":
		return DealStateFuture
	case "template":
		return DealStateTemplate
	}

	zap.L().Sugar().Warnf("invalid deal state (%s) -> use 'active' as default", state)

	return DealStateActive
}

func GetColorById(id int) string {
	switch id {
	case 1: // Elektronik & Technik
		return "#6898af"
	case 2: // Unterhaltung & Gaming
		return "#4774b2"
	case 3: // Lebensmittel & Haushalt
		return "#86b200"
	case 4: // Fashion, Schmuck & Lifestyle
		return "#b3396a"
	case 5: // Beauty, Wellness & Gesundheit
		return "#9059b3"
	case 6: // Family & Kids
		return "#02b0b2"
	case 7: // Home & Living
		return "#b2aba0"
	case 8: // Baumarkt & Garten
		return "#b28d4b"
	case 9: // Auto, Fahhrad & Motorrad
		return "#5c5e66"
	case 10: // Gastronomie, Bars & Cafes
		return "#b35a37"
	case 11: // Kultur & Freizeit
		return "#b3b100"
	case 12: // Sport & Outdoor
		return "#b22929"
	case 13: // Reisen, Hotels & Übernachtungen
		return "#3d484b"
	case 14: // Dienstleistungen & Finanzen
		return "#465c8e"
	case 15: // Floristik
		return "#60b262"
	case 16: // Sonstiges
		return "#b3b3b3"
	}

	zap.L().Sugar().Error("can't get color for deal with id: ", id)
	return "#ff00ff"
}

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
	PaymentState    *DealPaymentState `db:"payment_state"`
	StartInstantly  bool              `db:"start_instantly"`
	OwnEndDate      bool              `db:"own_end_date"`
}

func NewDeal() Deal {
	return Deal{
		Title:           "",
		Description:     "",
		CategoryId:      0,
		DurationInHours: 24,
		Start:           time.Now().Add(1 * time.Hour),
		IsTemplate:      false,
	}
}

func DealFromRequest(c echo.Context) (Deal, string) {
	title := c.FormValue("title")
	if len(title) == 0 {
		return Deal{}, "alert.enter_title"
	}

	description := c.FormValue("description")
	if len(description) == 0 {
		return Deal{}, "alert.enter_description"
	}

	categoryId, err := strconv.Atoi(c.FormValue("category"))
	if err != nil {
		return Deal{}, "alert.select_category"
	}

	isTemplate := c.FormValue("template") == "on"
	ownEndDate := c.FormValue("ownEndDate") == "on"
	startInstantly := c.FormValue("startInstantly") == "on"

	startDate := time.Now()

	if !startInstantly {
		startDate, err = time.Parse("2006-01-02T15:04", c.FormValue("startDate"))
		if err != nil {
			return Deal{}, "alert.provide_start_date"
		}
	}

	duration := 0
	if ownEndDate {
		endDate, err := time.Parse("2006-01-02", c.FormValue("endDate"))
		if err != nil {
			return Deal{}, "alert.provide_runtime_or_enddate"
		}
		duration = int(endDate.Sub(startDate.Truncate(24 * time.Hour)).Hours())
	} else {
		duration, err = strconv.Atoi(c.FormValue("duration"))
		if err != nil {
			return Deal{}, "alert.provide_runtime_or_enddate"
		}
		duration *= 24
	}

	if duration <= 0 {
		return Deal{}, "alert.start_before_end"
	}

	return Deal{
		Title:           title,
		Description:     description,
		CategoryId:      categoryId,
		Start:           startDate,
		DurationInHours: duration,
		IsTemplate:      isTemplate,
		OwnEndDate:      ownEndDate,
		StartInstantly:  startInstantly,
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

type DealHeader struct {
	ID         uuid.UUID
	DealerId   uuid.UUID `db:"dealer_id"`
	Title      string
	Username   string
	CategoryId int       `db:"category_id"`
	StartTime  time.Time `db:"start_time"`
	CanEdit    bool
}

type DealHeaders []DealHeader

func (deals DealHeaders) RotateByTime() DealHeaders {
	now := time.Now().Format("15:04")

	rotateIndex := 0
	for i, deal := range deals {
		dealStartTime := deal.StartTime.Format("15:04")
		if dealStartTime >= now {
			rotateIndex = i - 1
			break
		}
		rotateIndex = i
	}

	if rotateIndex >= 0 {
		return append(deals[rotateIndex:], deals[:rotateIndex]...)
	}

	return append(deals[len(deals)-1:], deals[:len(deals)-1]...)
}

type DealDetails struct {
	ID          uuid.UUID
	Title       string
	Description string
	Start       string
	End         string
}

type DealReport struct {
	Title  string
	Reason string
}

type DealStatistics struct {
	DealId          uuid.UUID `db:"deal_id"`
	DealerId        uuid.UUID `db:"dealer_id"`
	Title           string
	Start           time.Time
	DurationInHours int `db:"duration_in_hours"`
	LikeCount       int
	FavoriteCount   int
	ClickCount      int
}
