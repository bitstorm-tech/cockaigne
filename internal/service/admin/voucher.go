package adminservice

import (
	"github.com/bitstorm-tech/cockaigne/internal/model"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
)

func CreateVoucher(voucher model.Voucher) error {
	_, err := persistence.DB.Exec(
		"insert into vouchers (code, start, \"end\", discount_in_percent, is_active, multi_use, comment) values ($1, $2, $3, $4, $5, $6, $7)",
		voucher.Code,
		voucher.Start,
		voucher.End,
		voucher.DiscountInPercent,
		voucher.IsActive,
		voucher.MultiUse,
		voucher.Comment,
	)

	return err
}

func GetVouchers() ([]model.Voucher, error) {
	var vouchers []model.Voucher
	err := persistence.DB.Select(&vouchers, "select * from vouchers")
	if err != nil {
		return nil, err
	}

	return vouchers, nil
}
