package source

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/forbole/callisto/v4/types"
)

type Source interface {
	GetBalances(addresses []string, height int64) ([]types.AccountBalance, error)
	GetSupply(height int64) (sdk.Coins, error)

	// -- For hasura action --
	GetAccountBalance(address string, height int64) ([]sdk.Coin, error)
	GetDenomOwners(height int64, denom string) ([]*banktypes.DenomOwner, error)
}
