package database_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	dbtypes "github.com/forbole/callisto/v4/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveSupply() {
	// Save the data
	original := sdk.NewCoins(
		sdk.NewCoin("desmos", math.NewInt(10000)),
		sdk.NewCoin("uatom", math.NewInt(15)),
	)
	err := suite.database.SaveSupply(original, 10)
	suite.Require().NoError(err)

	// Verify the data
	expected := dbtypes.NewSupplyRow(dbtypes.NewDbCoins(original), 10)

	var rows []dbtypes.SupplyRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM supply`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "supply table should contain only one row")
	suite.Require().True(expected.Equals(rows[0]))

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with a lower height
	coins := sdk.NewCoins(
		sdk.NewCoin("desmos", math.NewInt(10000)),
		sdk.NewCoin("uatom", math.NewInt(15)),
	)
	err = suite.database.SaveSupply(coins, 9)
	suite.Require().NoError(err)

	// Verify the data
	rows = []dbtypes.SupplyRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM supply`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "supply table should contain only one row")
	suite.Require().True(expected.Equals(rows[0]))

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with same height
	coins = sdk.NewCoins(sdk.NewCoin("uakash", math.NewInt(10)))
	err = suite.database.SaveSupply(coins, 10)
	suite.Require().NoError(err)

	// Verify the data
	expected = dbtypes.NewSupplyRow(dbtypes.NewDbCoins(coins), 10)

	rows = []dbtypes.SupplyRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM supply`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "supply table should contain only one row")
	suite.Require().True(expected.Equals(rows[0]))

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with higher height
	coins = sdk.NewCoins(sdk.NewCoin("btc", math.NewInt(10)))
	err = suite.database.SaveSupply(coins, 20)
	suite.Require().NoError(err)

	// Verify the data
	expected = dbtypes.NewSupplyRow(dbtypes.NewDbCoins(coins), 20)

	rows = []dbtypes.SupplyRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM supply`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "supply table should contain only one row")
	suite.Require().True(expected.Equals(rows[0]))
}
