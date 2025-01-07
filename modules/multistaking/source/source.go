package source

import (
	multistakingtypes "github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

type Source interface {
	GetMultiStakingLocks(height int64) ([]*multistakingtypes.MultiStakingLock, error)
	GetMultiStakingUnlocks(height int64) ([]*multistakingtypes.MultiStakingUnlock, error)
	GetValidators(height int64, status string) ([]multistakingtypes.ValidatorInfo, error)
}
