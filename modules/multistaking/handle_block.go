package multistaking

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"

	juno "github.com/forbole/juno/v6/types"

	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	dbtypes "github.com/forbole/callisto/v4/database/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, res *tmctypes.ResultBlockResults, _ []*juno.Transaction, _ *tmctypes.ResultValidators,
) error {

	err := m.updateTxsByEvent(block.Block.Height, res.FinalizeBlockEvents)
	if err != nil {
		fmt.Printf("Error when updateTxsByEvent, error: %s", err)
	}
	return nil
}

func (m *Module) updateTxsByEvent(height int64, events []abci.Event) error {
	log.Debug().Str("module", "multistaking").Int64("height", height).
		Msg("updating txs by event")

	var msEvents []dbtypes.MSEvent
	for _, event := range events {
		switch event.Type {
		case stakingtypes.EventTypeDelegate:
			valAddr, _ := juno.FindAttributeByKey(event, stakingtypes.AttributeKeyValidator)
			delAddr, _ := juno.FindAttributeByKey(event, stakingtypes.AttributeKeyDelegator)
			amount, _ := juno.FindAttributeByKey(event, sdk.AttributeKeyAmount)
			msEvent, err := dbtypes.NewMSEvent("delegate", valAddr.Value, delAddr.Value, amount.Value)

			if err == nil {
				msEvents = append(msEvents, msEvent)
			}

		case stakingtypes.EventTypeUnbond:
			valAddr, _ := juno.FindAttributeByKey(event, stakingtypes.AttributeKeyValidator)
			delAddr, _ := juno.FindAttributeByKey(event, stakingtypes.AttributeKeyDelegator)
			amount, _ := juno.FindAttributeByKey(event, sdk.AttributeKeyAmount)
			msEvent, err := dbtypes.NewMSEvent("unbond", valAddr.Value, delAddr.Value, amount.Value)

			if err == nil {
				msEvents = append(msEvents, msEvent)
			}

		case stakingtypes.EventTypeCancelUnbondingDelegation:
			valAddr, _ := juno.FindAttributeByKey(event, stakingtypes.AttributeKeyValidator)
			delAddr, _ := juno.FindAttributeByKey(event, stakingtypes.AttributeKeyDelegator)
			amount, _ := juno.FindAttributeByKey(event, sdk.AttributeKeyAmount)
			msEvent, err := dbtypes.NewMSEvent("cancel_unbond", valAddr.Value, delAddr.Value, amount.Value)

			if err == nil {
				msEvents = append(msEvents, msEvent)
			}
		case stakingtypes.EventTypeCompleteRedelegation:
			valAddr1, _ := juno.FindAttributeByKey(event, stakingtypes.AttributeKeySrcValidator)
			valAddr2, _ := juno.FindAttributeByKey(event, stakingtypes.AttributeKeyDstValidator)
			delAddr, _ := juno.FindAttributeByKey(event, stakingtypes.AttributeKeyDelegator)
			m.UpdateLockAndUnlockInfo(height, delAddr.Value, valAddr1.Value)
			m.UpdateLockAndUnlockInfo(height, delAddr.Value, valAddr2.Value)

		case stakingtypes.EventTypeCompleteUnbonding:
			valAddr, _ := juno.FindAttributeByKey(event, stakingtypes.AttributeKeyValidator)
			delAddr, _ := juno.FindAttributeByKey(event, stakingtypes.AttributeKeyDelegator)
			m.CompleteUnbonding(height, delAddr.Value, valAddr.Value)
		}
	}

	if len(msEvents) == 0 {
		return nil
	}
	return m.db.SaveMSEvent(msEvents, height)
}
