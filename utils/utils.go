package utils

import (
	"context"
	"strconv"

	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"google.golang.org/grpc/metadata"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RemoveDuplicateValues removes the duplicated values from the given slice
func RemoveDuplicateValues(slice []string) []string {
	keys := make(map[string]bool)
	var list []string

	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// GetHeightRequestContext adds the height to the context for queries
func GetHeightRequestContext(context context.Context, height int64) context.Context {
	return metadata.AppendToOutgoingContext(
		context,
		grpctypes.GRPCBlockHeightHeader,
		strconv.FormatInt(height, 10),
	)
}

// ConvertValidatorPubKeyToBech32String converts the given pubKey to a Bech32 string
func ConvertValidatorPubKeyToBech32String(pubKey cryptotypes.PubKey) (string, error) {
	bech32Prefix := sdk.GetConfig().GetBech32ConsensusPubPrefix()
	return bech32.ConvertAndEncode(bech32Prefix, pubKey.Bytes())
}
