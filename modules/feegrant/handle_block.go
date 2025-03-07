package feegrant

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"

	feegranttypes "cosmossdk.io/x/feegrant"
	juno "github.com/forbole/juno/v6/types"

	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/callisto/v4/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, _ []*juno.Transaction, _ *tmctypes.ResultValidators,
) error {

	// Remove expired fee grant allowances
	err := m.removeExpiredFeeGrantAllowances(block.Block.Height, res.FinalizeBlockEvents)
	if err != nil {
		fmt.Printf("Error when removing expired fee grant allowance, error: %s", err)
	}
	return nil
}

// removeExpiredFeeGrantAllowances removes fee grant allowances in database that have expired
func (m *Module) removeExpiredFeeGrantAllowances(height int64, events []abci.Event) error {
	log.Debug().Str("module", "feegrant").Int64("height", height).
		Msg("updating expired fee grant allowances")

	events = juno.FindEventsByType(events, feegranttypes.EventTypeRevokeFeeGrant)

	for _, event := range events {
		granterAddress, err := juno.FindAttributeByKey(event, feegranttypes.AttributeKeyGranter)
		if err != nil {
			return fmt.Errorf("error while getting fee grant granter address: %s", err)
		}
		granteeAddress, err := juno.FindAttributeByKey(event, feegranttypes.AttributeKeyGrantee)
		if err != nil {
			return fmt.Errorf("error while getting fee grant grantee address: %s", err)
		}
		err = m.db.DeleteFeeGrantAllowance(types.NewGrantRemoval(granteeAddress.Value, granterAddress.Value, height))
		if err != nil {
			return fmt.Errorf("error while deleting fee grant allowance: %s", err)

		}
	}
	return nil

}
