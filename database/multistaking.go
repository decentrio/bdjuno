package database

import (
	"fmt"

	cosmossdk_io_math "cosmossdk.io/math"
	dbtypes "github.com/forbole/callisto/v4/database/types"

	"github.com/lib/pq"
	multistakingtypes "github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (db *Db) SaveMultiStakingLocks(height int64, multiStakingLocks []*multistakingtypes.MultiStakingLock) error {
	if len(multiStakingLocks) == 0 {
		return nil
	}

	query := `INSERT INTO ms_locks (staker_addr, val_addr, ms_lock, height) VALUES`

	var param []interface{}

	for i, msLock := range multiStakingLocks {
		vi := i * 4
		query += fmt.Sprintf("($%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4)
		mStakerAddr := msLock.LockID.MultiStakerAddr
		valAddr := msLock.LockID.ValAddr
		msCoin := dbtypes.NewMSCoin(msLock.LockedCoin)
		var mscoins dbtypes.MSCoins
		mscoins = append(mscoins, &msCoin)
		param = append(param, mStakerAddr, valAddr, pq.Array(mscoins), height)
	}

	query = query[:len(query)-1] // Remove trailing ","
	query += `
ON CONFLICT (staker_addr, val_addr) DO UPDATE 
	SET ms_lock = excluded.ms_lock,
		height = excluded.height
WHERE ms_locks.height <= excluded.height`

	_, err := db.SQL.Exec(query, param...)
	if err != nil {
		return fmt.Errorf("error while saving msLock: %s", err)
	}

	return nil
}

func (db *Db) SaveMultiStakingUnlocks(height int64, multiStakingUnlocks []*multistakingtypes.MultiStakingUnlock) error {
	if len(multiStakingUnlocks) == 0 {
		return nil
	}

	query := `INSERT INTO ms_unlocks (staker_addr, val_addr, unlock_entry, height) VALUES`

	var param []interface{}

	for i, msUnlock := range multiStakingUnlocks {
		vi := i * 4
		query += fmt.Sprintf("($%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4)
		mStakerAddr := msUnlock.UnlockID.MultiStakerAddr
		valAddr := msUnlock.UnlockID.ValAddr
		entries := msUnlock.Entries
		param = append(param, mStakerAddr, valAddr, pq.Array(dbtypes.NewUnlockEntries(entries)), height)
	}

	query = query[:len(query)-1] // Remove trailing ","
	query += `
ON CONFLICT (staker_addr, val_addr) DO UPDATE 
	SET unlock_entry = excluded.unlock_entry,
		height = excluded.height
WHERE ms_unlocks.height <= excluded.height`

	_, err := db.SQL.Exec(query, param...)
	if err != nil {
		return fmt.Errorf("error while saving msUnlock: %s", err)
	}

	return nil
}

func (db *Db) SaveUnbondingToken(height int64, multiStakingUnlocks []*multistakingtypes.MultiStakingUnlock) error {
	total := make(map[string]cosmossdk_io_math.Int)

	for _, msUnlock := range multiStakingUnlocks {
		entries := msUnlock.Entries
		for _, entry := range entries {
			denom := entry.UnlockingCoin.Denom
			amount := entry.UnlockingCoin.Amount
			if total[denom].IsNil() {
				total[denom] = amount
			} else {
				total[denom].Add(amount)
			}
		}
	}

	if len(total) == 0 {
		return nil
	}

	query := `INSERT INTO token_unbonding (denom, amount, height) VALUES`

	var param []interface{}

	i := 0
	for denom, amount := range total {
		vi := i * 3
		query += fmt.Sprintf("($%d,$%d,$%d),", vi+1, vi+2, vi+3)

		param = append(param, denom, amount.String(), height)
		i++
	}

	query = query[:len(query)-1] // Remove trailing ","
	query += `
ON CONFLICT (denom) DO UPDATE 
	SET amount = excluded.amount,
		height = excluded.height
WHERE token_unbonding.height <= excluded.height`

	_, err := db.SQL.Exec(query, param...)
	if err != nil {
		return fmt.Errorf("error while saving token_unbonding: %s", err)
	}

	return nil
}

func (db *Db) SaveBondedToken(height int64, multiStakingLocks []*multistakingtypes.MultiStakingLock) error {
	total := make(map[string]cosmossdk_io_math.Int)

	for _, msLock := range multiStakingLocks {
		denom := msLock.LockedCoin.Denom
		amount := msLock.LockedCoin.Amount
		if total[denom].IsNil() {
			total[denom] = amount
		} else {
			total[denom].Add(amount)
		}
	}

	if len(total) == 0 {
		return nil
	}

	query := `INSERT INTO token_bonded (denom, amount, height) VALUES`

	var param []interface{}

	i := 0
	for denom, amount := range total {
		vi := i * 3
		query += fmt.Sprintf("($%d,$%d,$%d),", vi+1, vi+2, vi+3)

		param = append(param, denom, amount.String(), height)
		i++
	}

	query = query[:len(query)-1] // Remove trailing ","
	query += `
ON CONFLICT (denom) DO UPDATE 
	SET amount = excluded.amount,
		height = excluded.height
WHERE token_bonded.height <= excluded.height`

	_, err := db.SQL.Exec(query, param...)
	if err != nil {
		return fmt.Errorf("error while saving token_bonded: %s", err)
	}

	return nil
}

func (db *Db) SaveValidatorDenom(height int64, validatorInfo []multistakingtypes.ValidatorInfo) error {
	if len(validatorInfo) == 0 {
		return nil
	}

	query := `INSERT INTO validator_denom (val_addr, denom, height) VALUES`

	var param []interface{}
	for i, info := range validatorInfo {
		vi := i * 3
		query += fmt.Sprintf("($%d,$%d,$%d),", vi+1, vi+2, vi+3)
		param = append(param, info.OperatorAddress, info.BondDenom, height)
	}

	query = query[:len(query)-1] // Remove trailing ","
	query += `
ON CONFLICT (val_addr) DO UPDATE 
	SET denom = excluded.denom,
    	height = excluded.height
WHERE validator_denom.height <= excluded.height`

	_, err := db.SQL.Exec(query, param...)
	if err != nil {
		return fmt.Errorf("error while saving ValidatorDenom: %s", err)
	}

	return nil
}

func (db *Db) SaveMSEvent(msEvents []dbtypes.MSEvent, height int64) error {
	if len(msEvents) == 0 {
		return nil
	}

	query := `INSERT INTO ms_event (height, name, val_addr, del_addr, amount) VALUES`

	var param []interface{}
	for i, msEvent := range msEvents {
		vi := i * 5
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4, vi+5)
		param = append(param, height, msEvent.Name, msEvent.ValAddr, msEvent.DelAddr, msEvent.Amount)
	}

	query = query[:len(query)-1] // Remove trailing ","

	_, err := db.SQL.Exec(query, param...)
	if err != nil {
		return fmt.Errorf("error while saving msEvents: %s", err)
	}

	return nil
}
