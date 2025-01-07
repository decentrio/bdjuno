package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v6/node/local"
	multistakingtypes "github.com/realio-tech/multi-staking-module/x/multi-staking/types"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/callisto/v4/modules/multistaking/source"
)

var (
	_ source.Source = &Source{}
)

// Source represents the implementation of the multistaking keeper that works on a local node
type Source struct {
	*local.Source
	qs multistakingtypes.QueryServer
}

// NewSource builds a new Source instance
func NewSource(source *local.Source, qs multistakingtypes.QueryServer) *Source {
	return &Source{
		Source: source,
		qs:     qs,
	}
}

func (s Source) GetMultiStakingLocks(height int64) ([]*multistakingtypes.MultiStakingLock, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var multiStakingLock []*multistakingtypes.MultiStakingLock
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.qs.MultiStakingLocks(
			sdk.WrapSDKContext(ctx),
			&multistakingtypes.QueryMultiStakingLocksRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 multiStakingLock at time
				},
			})
		if err != nil {
			return nil, fmt.Errorf("error while getting multiStakingLock: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		multiStakingLock = append(multiStakingLock, res.Locks...)
	}

	return multiStakingLock, nil
}

func (s Source) GetMultiStakingUnlocks(height int64) ([]*multistakingtypes.MultiStakingUnlock, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var multiStakingUnlock []*multistakingtypes.MultiStakingUnlock
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.qs.MultiStakingUnlocks(
			sdk.WrapSDKContext(ctx),
			&multistakingtypes.QueryMultiStakingUnlocksRequest{
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 multiStakingLock at time
				},
			})
		if err != nil {
			return nil, fmt.Errorf("error while getting multiStakingLock: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		multiStakingUnlock = append(multiStakingUnlock, res.Unlocks...)
	}

	return multiStakingUnlock, nil
}

func (s Source) GetValidators(height int64, status string) ([]multistakingtypes.ValidatorInfo, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	var validatorInfo []multistakingtypes.ValidatorInfo
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := s.qs.Validators(
			sdk.WrapSDKContext(ctx),
			&multistakingtypes.QueryValidatorsRequest{
				Status: status,
				Pagination: &query.PageRequest{
					Key:   nextKey,
					Limit: 100, // Query 100 validatorInfo at time
				},
			})
		if err != nil {
			return nil, fmt.Errorf("error while getting validatorInfo: %s", err)
		}

		nextKey = res.Pagination.NextKey
		stop = len(res.Pagination.NextKey) == 0
		validatorInfo = append(validatorInfo, res.Validators...)
	}

	return validatorInfo, nil
}
