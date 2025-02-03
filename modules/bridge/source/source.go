package source

import (
	"github.com/realiotech/realio-network/x/bridge/types"
)

type Source interface {
	RateLimits(height int64) ([]types.DenomAndRateLimit, error)
}
