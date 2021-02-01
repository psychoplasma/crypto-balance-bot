package cryptobot

import "math/big"

// CurrencyService represents API to fetch relavent info about account for the given currency
type CurrencyService interface {
	// GetAccountMovements fetches txs of the given address since the given block height(inclusive)
	// and converts it to account movements if there are any
	GetAccountMovements(address string, sinceBlockHeight uint64) (*AccountMovements, error)
	// GetLatestBlockHeight fetches the latest block number of the corresponding blockchain
	GetLatestBlockHeight() (uint64, error)
}

// Currency is a value object
type Currency struct {
	Symbol  string   `json:"symbol"`
	Decimal *big.Int `json:"decimal"`
}
