package bridge

import (
	"fmt"

	"github.com/rs/zerolog/log"
	juno "github.com/forbole/juno/v6/types"

	"github.com/forbole/callisto/v4/utils"

	bridgetypes "github.com/realiotech/realio-network/x/bridge/types"
)

var msgFilter = map[string]bool{
	"/realionetwork.bridge.v1.MsgBridgeIn":  true,
	"/realionetwork.bridge.v1.MsgBridgeOut": true,
}

// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ int, executedMsg juno.Message, tx *juno.Transaction) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(_ int, msg juno.Message, tx *juno.Transaction) error {
	if _, ok := msgFilter[msg.GetType()]; !ok {
		return nil
	}

	log.Debug().Str("module", "bridge").Str("hash", tx.TxHash).Uint64("height", tx.Height).Msg(fmt.Sprintf("handling bridge message %s", msg.GetType()))

	switch msg.GetType() {
	case "/realionetwork.bridge.v1.MsgBridgeIn":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &bridgetypes.MsgBridgeIn{})
		return m.db.SaveBridgeIn(tx.TxHash, cosmosMsg)

	case "/realionetwork.bridge.v1.MsgBridgeOut":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &bridgetypes.MsgBridgeOut{})
		return m.db.SaveBridgeOut(tx.TxHash, cosmosMsg)
	}
	return nil
}
