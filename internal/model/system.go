package model

import (
	"time"

	"github.com/google/uuid"
)

type Voucher struct {
	Code              string
	Comment           string
	Start             time.Time
	End               time.Time
	DiscountInPercent int  `db:"discount_in_percent"`
	IsActive          bool `db:"is_active"`
	MultiUse          bool `db:"multi_use"`
	Created           time.Time
}

type ActiveVoucher struct {
	AccountId         uuid.UUID `db:"account_id"`
	Activated         bool
	Code              string
	Start             time.Time
	End               time.Time
	DiscountInPercent int `db:"discount_in_percent"`
}

type RedeemedVoucher struct {
	AccountId   uuid.UUID `db:"account_id"`
	VoucherCode string    `db:"voucher_code"`
	RedeemedAt  time.Time `db:"redeemed_at"`
}
