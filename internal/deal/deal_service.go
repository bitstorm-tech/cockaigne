package deal

import (
	"fmt"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/gofiber/fiber/v2/log"
	"strings"
)

type State string

const (
	Active State = "active"
	Past   State = "past"
	Future State = "future"
)

func ToState(state string) State {
	switch strings.ToLower(state) {
	case "active":
		return Active
	case "past":
		return Past
	case "future":
		return Future
	}

	log.Warnf("invalid deal state (%s) -> use active as default", state)

	return Active
}

func GetCategories() []Category {
	var categories []Category
	err := persistence.DB.Select(&categories, "select * from categories where is_active = true")
	if err != nil {
		log.Errorf("can't get categories: %v", err)
	}

	return categories
}

func GetCategory(id int) (Category, error) {
	var category Category
	err := persistence.DB.Get(&category, "select * from categories where id = $1", id)
	if err != nil {
		return Category{}, fmt.Errorf("can't get category (id=%d): %v", id, err)
	}

	return category, nil
}

func SaveDeal(deal Deal) error {
	_, err := persistence.DB.Exec(
		"insert into deals (dealer_id, title, description, category_id, duration_in_hours, start, template) values ($1, $2, $3, $4, $5, $6, false)",
		deal.DealerId,
		deal.Title,
		deal.Description,
		deal.CategoryId,
		deal.DurationInHours,
		deal.Start,
	)

	if err != nil {
		return err
	}

	if deal.IsTemplate {
		_, err = persistence.DB.Exec(
			"insert into deals (dealer_id, title, description, category_id, duration_in_hours, start, template) values ($1, $2, $3, $4, $5, $6, true)",
			deal.DealerId,
			deal.Title,
			deal.Description,
			deal.CategoryId,
			deal.DurationInHours,
			deal.Start,
		)
	}

	return err
}

func GetDeal(id string) (Deal, error) {
	var deal Deal
	err := persistence.DB.Get(&deal, "select * from deals where id = $1", id)
	if err != nil {
		return Deal{}, fmt.Errorf("can't get deal from database: %v", err)
	}

	return deal, nil
}

func GetDealsFromView(state State, dealerId *string) ([]DealView, error) {
	if state != Future && state != Active && state != Past {
		return []DealView{}, fmt.Errorf("unknown deal state: %s", state)
	}

	statement := fmt.Sprintf("select *, st_x(location) || ',' || st_y(location) as location from %s_deals_view", state)

	if dealerId != nil {
		statement += fmt.Sprintf(" where dealer_id = '%s'", *dealerId)
	}

	var deals []DealView
	err := persistence.DB.Select(&deals, statement)

	if err != nil {
		return []DealView{}, fmt.Errorf("can't get active deals: %v", err)
	}

	return deals, nil
}

func GetTemplates(dealerId string) ([]Deal, error) {
	var templates []Deal
	err := persistence.DB.Select(&templates, "select * from deals where template = true and dealer_id = $1", dealerId)
	if err != nil {
		return []Deal{}, fmt.Errorf("can't get templates: %v", err)
	}

	return templates, nil
}
