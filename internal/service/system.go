package service

import (
	"errors"
	"fmt"

	"github.com/bitstorm-tech/cockaigne/internal/model"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
)

var ErrVoucherNotActive = errors.New("voucher currently not active")
var ErrVoucherAlreadyRedeemed = errors.New("voucher already redeemed")
var ErrVoucherCannotBeRedeemed = errors.New("voucher can't be redeemed")

func SaveContactMessage(accountId string, message string) error {
	if len(message) > 1000 {
		return fmt.Errorf("contact message greater then 1000 characters: %d", len(message))
	}

	_, err := persistence.DB.Exec("insert into contact_messages (account_id, message) values ($1, $2)", accountId, message)

	return err
}

func IsLastContactMessageYoungerThen5Minutes(accountId string) (bool, error) {
	messageYoungerThen5Minutes := true
	err := persistence.DB.Get(
		&messageYoungerThen5Minutes,
		"select exists (select * from contact_messages where account_id = $1 and created >= (now() - interval '5 minutes'))",
		accountId,
	)
	if err != nil {
		return true, err
	}

	return messageYoungerThen5Minutes, nil
}

func GetActiveVouchers(accountId string) ([]model.ActiveVoucher, error) {
	var activeVouchers []model.ActiveVoucher
	err := persistence.DB.Select(
		&activeVouchers,
		"select * from active_vouchers_view where account_id = $1",
		accountId,
	)
	if err != nil {
		return []model.ActiveVoucher{}, err
	}

	return activeVouchers, nil
}

func GetVoucher(voucherCode string) (model.Voucher, error) {
	var voucher model.Voucher
	err := persistence.DB.Get(
		&voucher,
		"select * from vouchers where code = $1",
		voucherCode,
	)

	return voucher, err
}

func RedeemVoucher(accountId string, voucherCode string) error {
	canBeRedeemed, err := CanVoucherBeRedeemed(voucherCode)
	if err != nil {
		return err
	}

	if !canBeRedeemed {
		return ErrVoucherCannotBeRedeemed
	}

	return SaveRedeemedVoucher(accountId, voucherCode)
}

func CanVoucherBeRedeemed(voucherCode string) (bool, error) {
	currentlyActive, err := VoucherCurrentlyActive(voucherCode)
	if err != nil {
		return false, err
	}
	if !currentlyActive {
		return false, ErrVoucherNotActive
	}

	alreadyRedeemed, err := VoucherAlreadyRedeemed(voucherCode)
	if err != nil || alreadyRedeemed {
		return false, err
	}
	if alreadyRedeemed {
		return false, ErrVoucherAlreadyRedeemed
	}

	return true, nil
}

func VoucherCurrentlyActive(voucherCode string) (bool, error) {
	var currentlyActive bool
	err := persistence.DB.Get(
		&currentlyActive,
		`select exists (select * from vouchers where code = $1 and now() between start and "end")`,
		voucherCode,
	)

	return currentlyActive, err
}

func VoucherAlreadyRedeemed(voucherCode string) (bool, error) {
	var multiUse bool
	err := persistence.DB.Get(
		&multiUse,
		"select multi_use from vouchers where code = $1",
		voucherCode,
	)
	if err != nil {
		return true, err
	}

	var redeemed bool
	err = persistence.DB.Get(
		&redeemed,
		"select exists (select * from redeemed_vouchers where code = $1)",
		voucherCode,
	)
	if err != nil {
		return true, err
	}

	return redeemed && !multiUse, nil
}

func SaveRedeemedVoucher(accountId string, voucherCode string) error {
	_, err := persistence.DB.Exec("insert into redeemed_vouchers (account_id, code) values ($1, $2)", accountId, voucherCode)

	return err
}
