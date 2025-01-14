package bank

import (
	"fmt"

	"github.com/forbole/callisto/v4/types"

	"github.com/rs/zerolog/log"
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

	var accountBalances []types.AccountBalance
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
			accountBalances = append(accountBalances, types.NewAccountBalance(addr, holder.Balance, block.Height))
		}
	}

	return m.db.SaveAccountBalances(accountBalances, block.Height)
}
