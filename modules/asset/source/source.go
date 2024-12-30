package source

import (
	assettypes "github.com/realiotech/realio-network/x/asset/types"
)

type Source interface {
	GetTokens(height int64) ([]assettypes.Token, error)
}
