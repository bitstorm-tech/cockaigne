package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID                   uuid.UUID
	Username             string
	Password             string
	Email                string
	Gender               sql.NullString
	Age                  sql.NullInt32
	IsDealer             bool `db:"is_dealer"`
	Street               sql.NullString
	HouseNumber          sql.NullString `db:"house_number"`
	City                 sql.NullString
	ZipCode              sql.NullInt32 `db:"zip"`
	Phone                sql.NullString
	TaxId                sql.NullString `db:"tax_id"`
	DefaultCategory      sql.NullInt32  `db:"default_category"`
	Location             sql.NullString
	Created              time.Time
	SearchRadiusInMeters int  `db:"search_radius_in_meters"`
	UseLocationService   bool `db:"use_location_service"`
}

type FavoriteCategory struct {
	AccountId  uuid.UUID
	CategoryId int
	CreatedAt  time.Time
}

type UpdateFilterRequest struct {
	SearchRadiusInMeters int
	FavoriteCategoryIds  []int
}
