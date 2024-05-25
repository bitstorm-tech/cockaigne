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
	Active               bool
	ActivationCode       sql.NullString `db:"activation_code"`
	ChangePasswordCode   sql.NullString `db:"change_password_code"`
	ChangeEmailCode      sql.NullString `db:"change_email_code"`
	NewEmail             sql.NullString `db:"new_email"`
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
	Language             string
}

type FavoriteCategory struct {
	AccountId  uuid.UUID
	CategoryId int
	CreatedAt  time.Time
}

type UpdateFilterRequest struct {
	SearchRadiusInMeters int   `form:"searchRadius"`
	FavoriteCategoryIds  []int `form:"favoriteCategoryIds"`
}
