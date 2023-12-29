package adminmodel

import (
	"database/sql"
	"errors"
	"time"

	"github.com/bitstorm-tech/cockaigne/internal/model"
	"go.uber.org/zap"
)

type CreateVoucherRequest struct {
	Code           string `form:"code"`
	Comment        string `form:"comment"`
	Start          string `form:"start"`
	End            string `form:"end"`
	DurationInDays int    `form:"durationInDays"`
	IsActive       string `form:"isActive"`
	MultiUse       string `form:"multiUse"`
	StartWhenEnter string `form:"startWhenEnter"`
}

func (c CreateVoucherRequest) ToVoucher() (model.Voucher, error) {
	var start sql.NullTime = sql.NullTime{Valid: false}
	if len(c.Start) > 0 {
		startDate, err := time.Parse("2006-01-02", c.Start)
		if err != nil {
			zap.L().Sugar().Error("can't parse start date from create voucher request: ", err)
		} else {
			start.Time = startDate
			start.Valid = true
		}
	}

	var end sql.NullTime = sql.NullTime{Valid: false}
	if len(c.End) > 0 {
		endDate, err := time.Parse("2006-01-02", c.End)
		if err != nil {
			zap.L().Sugar().Error("can't parse end date from create voucher request: ", err)
		} else {
			end.Time = endDate
			end.Valid = true
		}
	}

	if c.StartWhenEnter == "on" {
		start.Valid = false
	}

	if c.MultiUse == "on" && c.StartWhenEnter == "on" {
		return model.Voucher{}, errors.New("MultiUse and StartWhenEnter can't be set at the same time")
	}

	durationInDays := sql.NullInt32{Int32: int32(c.DurationInDays), Valid: c.DurationInDays > 0}

	if !durationInDays.Valid && !end.Valid {
		return model.Voucher{}, errors.New("either end date or duration has to be set")
	}

	return model.Voucher{
		Code:           c.Code,
		Comment:        c.Comment,
		Start:          start,
		End:            end,
		DurationInDays: durationInDays,
		IsActive:       c.IsActive == "on",
		MultiUse:       c.MultiUse == "on",
	}, nil
}
