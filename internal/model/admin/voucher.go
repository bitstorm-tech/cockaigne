package adminmodel

import (
	"time"

	"github.com/bitstorm-tech/cockaigne/internal/model"
)

type CreateVoucherRequest struct {
	Code              string `form:"code"`
	Comment           string `form:"comment"`
	Start             string `form:"start"`
	End               string `form:"end"`
	DiscountInPercent int    `form:"discountInPercent"`
	IsActive          string `form:"isActive"`
	MultiUse          string `form:"multiUse"`
}

func (c CreateVoucherRequest) ToVoucher() (model.Voucher, error) {
	start, err := time.Parse("2006-01-02", c.Start)
	if err != nil {
		return model.Voucher{}, err
	}

	end, err := time.Parse("2006-01-02", c.End)
	if err != nil {
		return model.Voucher{}, err
	}

	return model.Voucher{
		Code:              c.Code,
		Comment:           c.Comment,
		Start:             start,
		End:               end,
		DiscountInPercent: c.DiscountInPercent,
		IsActive:          c.IsActive == "on",
		MultiUse:          c.MultiUse == "on",
	}, nil
}
