package cryptobot

import "math/big"

// CurrencyService represents API to fetch relavent info about account for the given currency
type CurrencyService interface {
	GetBalance(addressDesc string) (*big.Int, error)
	GetTransactions(addressDesc string, since int) ([]*Transaction, error)
	CreateAddress(pubKey string) (string, error)
	DeriveAddressFromXPubKey(xPubKey string) ([]string, error)
	ValidateAddress(address string) error
	ValidatePubKey(pubKey string) error
}

// Currency is a value object
type Currency struct {
	Symbol  string `json:"symbol"`
	Decimal int    `json:"decimal"`
}
