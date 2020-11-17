package cryptobot

import "math/big"

// CurrencyAPI represents API to fetch relavent info about account for the given currency
type CurrencyAPI interface {
	GetBalance(addressDesc string) (*big.Int, error)
	GetTransactions(addressDesc string, since int) ([]string, error)
	CreateAddress(pubKey string) (string, error)
	ValidateAddress(address string) error
	ValidatePubKey(pubKey string) error
}

// Currency enumaration
type Currency struct {
	Decimal int
	Symbol  string
	API     CurrencyAPI
}
