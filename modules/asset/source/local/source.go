package local

import (
	"fmt"

	assettypes "github.com/realiotech/realio-network/x/asset/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v6/node/local"

	"github.com/forbole/callisto/v4/modules/asset/source"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the asset keeper that works on a local node
type Source struct {
	*local.Source
	q assettypes.QueryServer
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, bk assettypes.QueryServer) *Source {
	return &Source{
		Source: source,
		q:      bk,
	}
}

// GetSupply implements bankkeeper.Source
func (s Source) GetTokens(height int64) ([]assettypes.Token, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.q.Tokens(
		sdk.WrapSDKContext(ctx),
		&assettypes.QueryTokensRequest{})
	if err != nil {
		return nil, fmt.Errorf("error while getting tokens: %s", err)
	}

	return res.Tokens, nil
}
