package database

import (
	"fmt"

	dbtypes "github.com/forbole/callisto/v4/database/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/callisto/v4/types"
	"github.com/lib/pq"
)

// SaveSupply allows to save for the given height the given total amount of coins
func (db *Db) SaveSupply(coins sdk.Coins, height int64) error {
	query := `
INSERT INTO supply (coins, height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET coins = excluded.coins,
    	height = excluded.height
WHERE supply.height <= excluded.height`

	_, err := db.SQL.Exec(query, pq.Array(dbtypes.NewDbCoins(coins)), height)
	if err != nil {
		return fmt.Errorf("error while storing supply: %s", err)
	}

	return nil
}

func (db *Db) SaveTokenHolder(tokens map[string]int, height int64) error {
	if len(tokens) == 0 {
		return nil
	}

	query := `INSERT INTO token_holder (denom, num_holder, height) VALUES`

	var param []interface{}
	i := 0
	for denom, amount := range tokens {
		vi := i * 3
		query += fmt.Sprintf("($%d,$%d,$%d),", vi+1, vi+2, vi+3)
		param = append(param, denom, amount, height)
		i++
	}

	query = query[:len(query)-1] // Remove trailing ","
	query += `
ON CONFLICT (denom) DO UPDATE 
	SET num_holder = excluded.num_holder,
    	height = excluded.height
WHERE token_holder.height <= excluded.height`

	_, err := db.SQL.Exec(query, param...)
	if err != nil {
		return fmt.Errorf("error while saving token_holder: %s", err)
	}

	return nil
}

func (db *Db) SaveAccountBalances(accountBalances []types.AccountBalance, height int64) error {
	if len(accountBalances) == 0 {
		return nil
	}

	query := `INSERT INTO balance (address, balances, height) VALUES`

	var param []interface{}
	for i, accountBalance := range accountBalances {
		vi := i * 3
		query += fmt.Sprintf("($%d,$%d,$%d),", vi+1, vi+2, vi+3)
		param = append(param, accountBalance.Address, pq.Array(dbtypes.NewDbCoins(accountBalance.Balance)), height)
	}

	query = query[:len(query)-1] // Remove trailing ","
	query += `
ON CONFLICT (address) DO UPDATE 
	SET balances = excluded.balances,
    	height = excluded.height
WHERE balance.height <= excluded.height`

	_, err := db.SQL.Exec(query, param...)
	if err != nil {
		return fmt.Errorf("error while saving AccountBalances: %s", err)
	}

	return nil
}
