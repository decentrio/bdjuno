package remote

import (
	"fmt"

	"github.com/forbole/juno/v6/node/remote"

	assetsource "github.com/forbole/callisto/v4/modules/asset/source"
	assettypes "github.com/realiotech/realio-network/x/asset/types"
)

var (
	_ assetsource.Source = &Source{}
)

type Source struct {
	*remote.Source
	assetClient assettypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, assetClient assettypes.QueryClient) *Source {
	return &Source{
		Source:      source,
		assetClient: assetClient,
	}
}

// GetSupply implements assetsource.Source
func (s Source) GetTokens(height int64) ([]assettypes.Token, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	res, err := s.assetClient.Tokens(
		ctx,
		&assettypes.QueryTokensRequest{})
	if err != nil {
		return nil, fmt.Errorf("error while getting total supply: %s", err)
	}

	return res.Tokens, nil
}
