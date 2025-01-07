package asset

import (
	"fmt"
	"strings"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/callisto/v4/types"
	assettypes "github.com/realiotech/realio-network/x/asset/types"
)

// RegisterPeriodicOperations implements modules.Module
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "asset").Msg("setting up periodic tasks")

	err := m.InitAllTokens()
	if err != nil {
		return fmt.Errorf("error while initAllTokens: %s", err)
	}

	return nil
}

// InitAllTokens init the supply of all the tokens
func (m *Module) InitAllTokens() error {
	log.Trace().Str("module", "bank").Str("operation", "total supply").
		Msg("updating total supply")

	block, err := m.db.GetLastBlockHeightAndTimestamp()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	tokens, err := m.source.GetTokens(block.Height)
	if err != nil {
		return err
	}

	return m.updateAllTokens(block.Height, tokens)
}

// updateTokenHolder updates num holder of all the tokens
func (m *Module) updateAllTokens(height int64, tokens []assettypes.Token) error {
	log.Trace().Str("module", "asset").Str("operation", "token").
		Msg("updating token unit")

	for _, token := range tokens {
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
