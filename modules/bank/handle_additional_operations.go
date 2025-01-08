package bank

import (
	"fmt"

	"github.com/forbole/callisto/v4/types"

	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RunAdditionalOperations implements modules.AdditionalOperationsModule
func (m *Module) RunAdditionalOperations() error {

	return m.initAccountBalances()
}

func (m *Module) initAccountBalances() error {
	log.Trace().Str("module", "bank").Str("operation", "account balance").
		Msg("init account balance")

	block, err := m.db.GetLastBlockHeightAndTimestamp()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	tokens, err := m.keeper.GetSupply(block.Height)
	if err != nil {
		return err
	}

	totalBalance := make(map[string]sdk.Coins)
	for _, tokenUnit := range tokens {
		denom := tokenUnit.Denom
		holders, err := m.keeper.GetDenomOwners(block.Height, denom)
		if err != nil {
			return fmt.Errorf("error while get denom holder: %s", err)
		}

		for _, holder := range holders {
			addr := holder.GetAddress()
			if addr == "" {
				continue
			}

			if _, exists := totalBalance[addr]; !exists {
				totalBalance[addr] = sdk.Coins{}
			}
			totalBalance[addr] = totalBalance[addr].Add(holder.Balance)
		}
	}

	if len(totalBalance) == 0 {
		return nil
	}

	var accountBalances []types.AccountBalance
	for addr, balance := range totalBalance {
		accountBalances = append(accountBalances, types.NewAccountBalance(
			addr,
			balance,
			block.Height,
		))
	}

	return m.db.SaveAccountBalances(accountBalances, block.Height)
}
