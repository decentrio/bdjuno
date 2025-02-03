package bridge

import (
	"fmt"

	juno "github.com/forbole/juno/v6/types"

	"github.com/rs/zerolog/log"

	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
)

// HandleBlock implements modules.Module
func (m *Module) HandleBlock(
	b *tmctypes.ResultBlock, _ *tmctypes.ResultBlockResults, _ []*juno.Transaction, _ *tmctypes.ResultValidators,
) error {
	err := m.updateRateLimits(b)
	if err != nil {
		log.Error().Str("module", "bridge").Int64("height", b.Block.Height).
			Err(err).Msg("error while updating rate limits")
	}

	return nil
}

func (m *Module) updateRateLimits(block *tmctypes.ResultBlock) error {
	log.Trace().Str("module", "bridge").Int64("height", block.Block.Height).
		Msg("updating rate limits")

	rates, err := m.source.RateLimits(block.Block.Height)
	if err != nil {
		return fmt.Errorf("error while getting rate limits: %s", err)
	}

	return m.db.SaveRates(block.Block.Height, rates)
}
