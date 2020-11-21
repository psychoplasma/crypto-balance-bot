package coin

import (
	"errors"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

var ErrUnknownCurrency = errors.New("unknown currency")

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

// Factory gets the corresponding coin API implementation
func Factory(c domain.Currency) (domain.CurrencyService, error) {
	switch c {
	case BTC:
		return &BitcoinAPI{}, nil
	case ETH:
		return &EthereumAPI{}, nil
	default:
		return nil, ErrUnknownCurrency
	}
}
