package bank

import (
	"fmt"

	"strings"

	"github.com/forbole/callisto/v4/types"
	"github.com/forbole/callisto/v4/utils"
	juno "github.com/forbole/juno/v6/types"
	assettypes "github.com/realiotech/realio-network/x/asset/types"
	"github.com/rs/zerolog/log"
)

var msgFilter = map[string]bool{
	"/realionetwork.asset.v1.MsgCreateToken":  true,
	"/realionetwork.asset.v1.Msg/CreateToken": true,
}

// HandleMsg implements MessageModule
// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ int, executedMsg juno.Message, tx *juno.Transaction) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(_ int, msg juno.Message, tx *juno.Transaction) error {
	if _, ok := msgFilter[msg.GetType()]; !ok {
		return nil
	}

	log.Debug().Str("module", "asset").Str("hash", tx.TxHash).Uint64("height", tx.Height).Msg(fmt.Sprintf("handling create token message %s", msg.GetType()))

	switch msg.GetType() {
	case "/realionetwork.asset.v1.Msg/CreateToken":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &assettypes.MsgCreateToken{})
		return m.handleMsgCreateToken(cosmosMsg)

	case "/realionetwork.asset.v1.MsgCreateToken":
		cosmosMsg := utils.UnpackMessage(m.cdc, msg.GetBytes(), &assettypes.MsgCreateToken{})
		return m.handleMsgCreateToken(cosmosMsg)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// handleMsgCreateToken handles properly a MsgCreateToken instance by
// saving into the database all the data associated to such NewToken
func (m *Module) handleMsgCreateToken(msg *assettypes.MsgCreateToken) error {
	symbol := msg.Symbol
	if symbol == "" {
		return nil
	}

	lowerCaseSymbol := strings.ToLower(msg.Symbol)
	baseDenom := fmt.Sprintf("a%s", lowerCaseSymbol)
	tokenUnit := types.NewTokenUnit(baseDenom, 18, nil, "")
	token := types.NewToken(symbol, []types.TokenUnit{tokenUnit})
	return m.db.SaveToken(token)
}
