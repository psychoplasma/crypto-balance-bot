package cryptobot

// CurrencyService represents API to fetch relavent info about account for the given currency
type CurrencyService interface {
	// GetAddressTxs fetches txs of the given address since the given block height(exclusive)
	GetAddressTxs(address string, sinceBlockHeight int) ([]*AccountMovement, error)
}

// Currency is a value object
type Currency struct {
	Symbol  string `json:"symbol"`
	Decimal int    `json:"decimal"`
}
