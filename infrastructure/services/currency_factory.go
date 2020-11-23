package services

import (
	"errors"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchaindotcom"
)

var errUnknownCurrency = errors.New("unknown currency")

var (
	BTC = domain.Currency{
		Decimal: 8,
		Symbol:  "btc",
	}
	ETH = domain.Currency{
		Decimal: 16,
		Symbol:  "eth",
	}
)

// CurrencyFactory gets the corresponding currency API implementation
func CurrencyFactory(c domain.Currency) (domain.CurrencyService, error) {
	switch c {
	case BTC:
		return &blockchaindotcom.BitcoinAPI{
			T: blockchaindotcom.BitcoinTranslator{},
		}, nil
	default:
		return nil, errUnknownCurrency
	}
}

// Translator interface declares translation from outside service output to currency service inputs
type Translator interface {
	ToAccountMovements(address string, addressData interface{}) []*domain.AccountMovement
}
