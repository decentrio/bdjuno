package local

import (
	"fmt"

	"github.com/forbole/juno/v6/node/local"
	bridgetypes "github.com/realiotech/realio-network/x/bridge/types"

	"github.com/forbole/callisto/v4/modules/bridge/source"
)

var (
	_ source.Source = &Source{}
)

type Source struct {
	*local.Source
	qs bridgetypes.QueryServer
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, qs bridgetypes.QueryServer) *Source {
	return &Source{
		Source: source,
		qs:     qs,
	}
}

func (s Source) RateLimits(height int64) ([]bridgetypes.DenomAndRateLimit, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.qs.RateLimits(ctx, &bridgetypes.QueryRateLimitsRequest{})
	if err != nil {
		return []bridgetypes.DenomAndRateLimit{}, err
	}
	return res.Ratelimits, nil
}
