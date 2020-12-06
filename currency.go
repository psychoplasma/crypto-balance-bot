package cryptobot

import "math/big"

// CurrencyService represents API to fetch relavent info about account for the given currency
type CurrencyService interface {
	// GetTxsOfAddress fetches txs of the given address since the given block height(exclusive)
	GetTxsOfAddress(address string, sinceBlockHeight int) (*AccountMovements, error)
}

// Currency is a value object
type Currency struct {
	Symbol  string   `json:"symbol"`
	Decimal *big.Int `json:"decimal"`
}
