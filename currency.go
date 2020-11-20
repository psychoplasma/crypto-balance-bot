package cryptobot

import "math/big"

// CurrencyService represents API to fetch relavent info about account for the given currency
type CurrencyService interface {
	// GetBalance fetches balance of the given address
	GetBalance(address string) (*big.Int, error)
	// GetTransactions fetches transaction of the given address starting from the given index
	GetTransactions(address string, index int) ([]*Transaction, error)
	// ValidateAddress checks whether or not the given address is valid
	// and returns an error in case of invalid address
	ValidateAddress(address string) error
}

// Currency is a value object
type Currency struct {
	Symbol  string `json:"symbol"`
	Decimal int    `json:"decimal"`
}
