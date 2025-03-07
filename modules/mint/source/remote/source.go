package remote

import (
	"cosmossdk.io/math"
	"github.com/forbole/juno/v6/node/remote"
	minttypes "github.com/realiotech/realio-network/x/mint/types"

	mintsource "github.com/forbole/callisto/v4/modules/mint/source"
)

var (
	_ mintsource.Source = &Source{}
)

// Source implements mintsource.Source using a remote node
type Source struct {
	*remote.Source
	querier minttypes.QueryClient
}

// NewSource returns a new Source instance
func NewSource(source *remote.Source, querier minttypes.QueryClient) *Source {
	return &Source{
		Source:  source,
		querier: querier,
	}
}

// GetInflation implements mintsource.Source
func (s Source) GetInflation(height int64) (math.LegacyDec, error) {
	res, err := s.querier.Inflation(remote.GetHeightRequestContext(s.Ctx, height), &minttypes.QueryInflationRequest{})
	if err != nil {
		return math.LegacyDec{}, err
	}

	return res.Inflation, nil
}

// Params implements mintsource.Source
func (s Source) Params(height int64) (minttypes.Params, error) {
	res, err := s.querier.Params(remote.GetHeightRequestContext(s.Ctx, height), &minttypes.QueryParamsRequest{})
	if err != nil {
		return minttypes.Params{}, nil
	}

	return res.Params, nil
}
