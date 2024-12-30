package asset

import (
	"fmt"
	"strings"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/callisto/v4/modules/utils"
	"github.com/forbole/callisto/v4/types"
	assettypes "github.com/realiotech/realio-network/x/asset/types"
)

// RegisterPeriodicOperations implements modules.Module
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "asset").Msg("setting up periodic tasks")
	fmt.Println("vlin0")
	if _, err := scheduler.Every(10).Minutes().Do(func() {
		utils.WatchMethod(m.UpdateTokens)
	}); err != nil {
		return fmt.Errorf("error while setting up asset token periodic operation: %s", err)
	}

	return nil
}

// UpdateSupply updates the supply of all the tokens
func (m *Module) UpdateTokens() error {
	log.Trace().Str("module", "bank").Str("operation", "total supply").
		Msg("updating total supply")

	block, err := m.db.GetLastBlockHeightAndTimestamp()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}
	fmt.Println("vlin1")

	tokens, err := m.source.GetTokens(block.Height)
	if err != nil {
		return err
	}
	fmt.Println("vlin2", len(tokens))

	return m.updateAllTokens(block.Height, tokens)
}

// updateTokenHolder updates num holder of all the tokens
func (m *Module) updateAllTokens(height int64, tokens []assettypes.Token) error {
	log.Trace().Str("module", "asset").Str("operation", "token").
		Msg("updating token unit")

	for _, token := range tokens {
		fmt.Println("vlin3", token.Symbol)
		lowerCaseSymbol := strings.ToLower(token.Symbol)
		baseDenom := fmt.Sprintf("a%s", lowerCaseSymbol)
		tokenUnit := types.NewTokenUnit(baseDenom, 18, nil, "")
		token := types.NewToken(token.Symbol, []types.TokenUnit{tokenUnit})
		err := m.db.SaveToken(token)
		if err != nil {
			return fmt.Errorf("error while save token2 : %s", err)
		}
	}

	return nil
}
