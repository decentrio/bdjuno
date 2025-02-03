package types

type BridgeIn struct {
	Hash     string `db:"hash"`
	Amount   string `db:"amount"`
	Denom    string `db:"denom"`
	Receiver string `db:"receiver"`
}

type BridgeOut struct {
	Hash   string `db:"hash"`
	Amount string `db:"amount"`
	Denom  string `db:"denom"`
	Sender string `db:"sender"`
}
