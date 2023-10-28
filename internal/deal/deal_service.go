package deal

import (
	"fmt"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/gofiber/fiber/v2/log"
)

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
	err := persistence.DB.Get(category, "select * from categories where id = $1", id)
	if err != nil {
		return Category{}, fmt.Errorf("can't get category (id=%d): %v", id, err)
	}

	return category, nil
}

func SaveDeal(deal Deal) error {
	_, err := persistence.DB.Exec(
		"insert into deals (dealer_id, title, description, category_id, duration_in_hours, start) values ($1, $2, $3, $4, $5, $6)",
		deal.DealerId,
		deal.Title,
		deal.Description,
		deal.CategoryId,
		deal.DurationInHours,
		deal.Start,
	)

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

func GetActiveDeals() ([]ActiveDeal, error) {
	var deals []ActiveDeal
	err := persistence.DB.Select(&deals, "select *, st_x(location) || ',' || st_y(location) as location from active_deals_view")

	if err != nil {
		return []ActiveDeal{}, fmt.Errorf("can't get active deals: %v", err)
	}

	return deals, nil
}

func GetDealsOfDealer(dealerId string) ([]Deal, error) {
	var deals []Deal
	err := persistence.DB.Select(&deals, "select * from deals where dealer_id = $1", dealerId)
	if err != nil {
		return []Deal{}, fmt.Errorf("can't get deals of dealer (id=%s): %v", dealerId, err)
	}

	return deals, nil
}
