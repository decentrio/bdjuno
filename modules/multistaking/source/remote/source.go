package remote

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"
	multistakingsource "github.com/forbole/callisto/v4/modules/multistaking/source"
	"github.com/forbole/juno/v6/node/remote"
	multistakingtypes "github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

var (
	_ multistakingsource.Source = &Source{}
)

type Source struct {
	*remote.Source
	msClient multistakingtypes.QueryClient
}

// NewSource builds a new Source instance
func NewSource(source *remote.Source, msClient multistakingtypes.QueryClient) *Source {
	return &Source{
		Source:   source,
		msClient: msClient,
	}
}

func (s Source) GetMultiStakingLocks(height int64) ([]*multistakingtypes.MultiStakingLock, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var multiStakingLock []*multistakingtypes.MultiStakingLock
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.msClient.MultiStakingLocks(
			ctx,
			&multistakingtypes.QueryMultiStakingLocksRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 multiStakingLock at time
				},
			})
		if err != nil {
			return nil, fmt.Errorf("error while getting total supply: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		multiStakingLock = append(multiStakingLock, res.Locks...)
	}

	return multiStakingLock, nil
}

func (s Source) GetMultiStakingUnlocks(height int64) ([]*multistakingtypes.MultiStakingUnlock, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var multiStakingUnlock []*multistakingtypes.MultiStakingUnlock
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.msClient.MultiStakingUnlocks(
			ctx,
			&multistakingtypes.QueryMultiStakingUnlocksRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 multiStakingLock at time
				},
			})
		if err != nil {
			return nil, fmt.Errorf("error while getting total supply: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		multiStakingUnlock = append(multiStakingUnlock, res.Unlocks...)
	}

	return multiStakingUnlock, nil
}

func (s Source) GetValidators(height int64, status string) ([]multistakingtypes.ValidatorInfo, error) {
	ctx := remote.GetHeightRequestContext(s.Ctx, height)

	var validatorInfo []multistakingtypes.ValidatorInfo
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.msClient.Validators(
			ctx,
			&multistakingtypes.QueryValidatorsRequest{
				Status: status,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 validatorInfo at time
				},
			})
		if err != nil {
			return nil, fmt.Errorf("error while getting total validatorInfo: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		validatorInfo = append(validatorInfo, res.Validators...)
	}

	return validatorInfo, nil
}
