package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Voucher struct {
	Code           string
	Comment        string
	Start          sql.NullString
	End            sql.NullString
	DurationInDays int  `db:"duration_in_days"`
	IsActive       bool `db:"is_active"`
	MultiUse       bool `db:"multi_use"`
}

type ActiveVoucher struct {
	AccountId      uuid.UUID `db:"account_id"`
	Activated      bool
	Code           string
	Start          time.Time
	End            time.Time
	DurationInDays int `db:"duration_in_days"`
}
