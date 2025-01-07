package asset

import (
	"github.com/cosmos/cosmos-sdk/codec"
	assetsource "github.com/forbole/callisto/v4/modules/asset/source"
	"github.com/forbole/juno/v6/modules"

	"github.com/forbole/callisto/v4/database"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
	_ modules.MessageModule            = &Module{}
	_ modules.AuthzMessageModule       = &Module{}
)

// Module represents the x/staking module
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	source assetsource.Source
}

// NewModule returns a new Module instance
func NewModule(
	source assetsource.Source, cdc codec.Codec, db *database.Db,
) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		source: source,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "asset"
}
