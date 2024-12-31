package types

import (
	"database/sql/driver"
	"fmt"
	"strings"

	multistakingtypes "github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

type MSCoin struct {
	Denom      string
	Amount     string
	BondWeight string
}

type UnlockEntry struct {
	CreationHeight string
	Denom          string
	Amount         string
	BondWeight     string
}

func NewMSCoin(coin multistakingtypes.MultiStakingCoin) MSCoin {
	return MSCoin{
		Denom:      coin.Denom,
		Amount:     coin.Amount.String(),
		BondWeight: coin.BondWeight.String(),
	}
}

// Value implements driver.Valuer
func (coin *MSCoin) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%s,%s)", coin.Denom, coin.Amount, coin.BondWeight), nil
}

// Scan implements sql.Scanner
func (coin *MSCoin) Scan(src interface{}) error {
	strValue := string(src.([]byte))
	strValue = strings.ReplaceAll(strValue, `"`, "")
	strValue = strings.ReplaceAll(strValue, "{", "")
	strValue = strings.ReplaceAll(strValue, "}", "")
	strValue = strings.ReplaceAll(strValue, "(", "")
	strValue = strings.ReplaceAll(strValue, ")", "")

	values := strings.Split(strValue, ",")

	*coin = MSCoin{Denom: values[0], Amount: values[1], BondWeight: values[2]}
	return nil
}

type MSCoins []*MSCoin

// Scan implements sql.Scanner
func (coins *MSCoins) Scan(src interface{}) error {
	strValue := string(src.([]byte))
	strValue = strings.ReplaceAll(strValue, `"`, "")
	strValue = strings.ReplaceAll(strValue, "{", "")
	strValue = strings.ReplaceAll(strValue, "}", "")
	strValue = strings.ReplaceAll(strValue, "),(", ") (")
	strValue = strings.ReplaceAll(strValue, "(", "")
	strValue = strings.ReplaceAll(strValue, ")", "")

	values := RemoveEmpty(strings.Split(strValue, " "))

	coinsV := make(MSCoins, len(values))
	for index, value := range values {
		v := strings.Split(value, ",") // Split the values

		coin := MSCoin{Denom: v[0], Amount: v[1], BondWeight: v[2]}
		coinsV[index] = &coin
	}

	*coins = coinsV
	return nil
}

func NewUnlockEntry(entry multistakingtypes.UnlockEntry) UnlockEntry {
	return UnlockEntry{
		CreationHeight: fmt.Sprintf("%d", entry.GetCreationHeight()),
		Denom:          entry.UnlockingCoin.Denom,
		Amount:         entry.UnlockingCoin.Amount.String(),
		BondWeight:     entry.UnlockingCoin.BondWeight.String(),
	}
}

// Value implements driver.Valuer
func (entry *UnlockEntry) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%s,%s,%s)", entry.CreationHeight, entry.Denom, entry.Amount, entry.BondWeight), nil
}

// Scan implements sql.Scanner
func (entry *UnlockEntry) Scan(src interface{}) error {
	strValue := string(src.([]byte))
	strValue = strings.ReplaceAll(strValue, `"`, "")
	strValue = strings.ReplaceAll(strValue, "{", "")
	strValue = strings.ReplaceAll(strValue, "}", "")
	strValue = strings.ReplaceAll(strValue, "(", "")
	strValue = strings.ReplaceAll(strValue, ")", "")

	values := strings.Split(strValue, ",")
	*entry = UnlockEntry{CreationHeight: values[0], Denom: values[1], Amount: values[2], BondWeight: values[3]}
	return nil
}

type UnlockEntries []*UnlockEntry

func NewUnlockEntries(entries []multistakingtypes.UnlockEntry) UnlockEntries {
	unlockEntries := make([]*UnlockEntry, 0)
	for _, entry := range entries {
		unlockEntry := NewUnlockEntry(entry)
		unlockEntries = append(unlockEntries, &unlockEntry)
	}
	return unlockEntries
}

type MSEvent struct {
	Name    string
	ValAddr string
	DelAddr string
	Amount  string
}

func NewMSEvent(ValAddr string, DelAddr string, Amount string, Name string) (MSEvent, error) {
	if ValAddr != "" && DelAddr != "" && Amount != "" && Name != "" {
		return MSEvent{
			Name:    Name,
			ValAddr: ValAddr,
			DelAddr: DelAddr,
			Amount:  Amount,
		}, nil
	}

	return MSEvent{}, fmt.Errorf("error")
}

// Scan implements sql.Scanner
func (coins *UnlockEntries) Scan(src interface{}) error {
	strValue := string(src.([]byte))
	strValue = strings.ReplaceAll(strValue, `"`, "")
	strValue = strings.ReplaceAll(strValue, "{", "")
	strValue = strings.ReplaceAll(strValue, "}", "")
	strValue = strings.ReplaceAll(strValue, "),(", ") (")
	strValue = strings.ReplaceAll(strValue, "(", "")
	strValue = strings.ReplaceAll(strValue, ")", "")

	values := RemoveEmpty(strings.Split(strValue, " "))

	coinsV := make(UnlockEntries, len(values))
	for index, value := range values {
		v := strings.Split(value, ",") // Split the values

		coin := UnlockEntry{CreationHeight: v[0], Denom: v[1], Amount: v[2], BondWeight: v[3]}
		coinsV[index] = &coin
	}

	*coins = coinsV
	return nil
}
