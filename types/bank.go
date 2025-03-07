package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountBalance represents the balance of an account at a given height
type AccountBalance struct {
	Address string
	Denom   string
	Amount  string
	Height  int64
}

// NewAccountBalance allows to build a new AccountBalance instance
func NewAccountBalance(address string, balance sdk.Coin, height int64) AccountBalance {
	return AccountBalance{
		Address: address,
		Denom:   balance.Denom,
		Amount:  balance.Amount.String(),
		Height:  height,
	}
}
