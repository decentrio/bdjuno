package multistaking

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/forbole/callisto/v4/modules/utils"
)

// RegisterPeriodicOperations implements modules.Module
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "multistaking").Msg("setting up periodic tasks")

	if _, err := scheduler.Every(5).Hour().Do(func() {
		utils.WatchMethod(m.UpdateMultiStaking)
	}); err != nil {
		return fmt.Errorf("error while setting up multistaking token periodic operation: %s", err)
	}

	return nil
}

func (m *Module) UpdateMultiStaking() error {
	log.Trace().Str("module", "multistaking").Str("operation", "multistaking lock").
		Msg("updating multistaking lock")

	block, err := m.db.GetLastBlockHeightAndTimestamp()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	err = m.UpdateMultiStakingLocks(block.Height)
	if err != nil {
		return fmt.Errorf("error while update UpdateMultiStakingLocks: %s", err)
	}

	err = m.UpdateMultiStakingUnlocks(block.Height)
	if err != nil {
		return fmt.Errorf("error while update UpdateMultiStakingUnlocks: %s", err)
	}

	err = m.UpdateValidatorInfo(block.Height)
	if err != nil {
		return fmt.Errorf("error while update UpdateValidatorInfo: %s", err)
	}

	return nil
}

func (m *Module) UpdateMultiStakingLocks(height int64) error {
	log.Trace().Str("module", "multistaking").Str("operation", "multistaking lock").
		Msg("updating multistaking lock")

	multiStakingLocks, err := m.source.GetMultiStakingLocks(height)
	if err != nil {
		return err
	}

	err = m.db.SaveBondedToken(height, multiStakingLocks)
	if err != nil {
		return err
	}

	return m.db.SaveMultiStakingLocks(height, multiStakingLocks)
}

func (m *Module) UpdateMultiStakingUnlocks(height int64) error {
	log.Trace().Str("module", "multistaking").Str("operation", "multistaking unlock").
		Msg("updating multistaking unlock")

	multiStakingUnlocks, err := m.source.GetMultiStakingUnlocks(height)
	if err != nil {
		return err
	}

	err = m.db.SaveUnbondingToken(height, multiStakingUnlocks)
	if err != nil {
		return err
	}

	return m.db.SaveMultiStakingUnlocks(height, multiStakingUnlocks)
}

func (m *Module) UpdateValidatorInfo(height int64) error {
	log.Trace().Str("module", "multistaking").Str("operation", "validator info").
		Msg("updating validator info")

	validatorInfo, err := m.source.GetValidators(height, "")
	if err != nil {
		return err
	}

	return m.db.SaveValidatorDenom(height, validatorInfo)
}
