package remote

import (
	"fmt"

	bridgesource "github.com/forbole/callisto/v4/modules/bridge/source"
	"github.com/forbole/juno/v6/node/remote"
	bridgetypes "github.com/realiotech/realio-network/x/bridge/types"
)

var (
	_ bridgesource.Source = &Source{}
)

type Source struct {
	*remote.Source
	msClient bridgetypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, msClient bridgetypes.QueryClient) *Source {
	return &Source{
		Source:   source,
		msClient: msClient,
	}
}

func (s Source) RateLimits(height int64) ([]bridgetypes.DenomAndRateLimit, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	res, err := s.msClient.RateLimits(ctx, &bridgetypes.QueryRateLimitsRequest{})
	if err != nil {
		return nil, fmt.Errorf("error while getting bridge rate limit: %s", err)
	}

	return res.Ratelimits, nil
}
