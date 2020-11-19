package application

import (
	domain "github.com/psychoplasma/crypto-balance-bot"
)

// AvailableCurrencies contains the implemented currencies
// TODO: Load these through a configuration file
var availableCurrencies = map[string]*domain.Currency{
	"btc": {
		Decimal: 8,
		Symbol:  "btc",
	},
	"eth": {
		Decimal: 16,
		Symbol:  "eth",
	},
}

// CurrencyService exposes application services for currency entity
type CurrencyService struct {
}

// NewCurrencyService factory function
func NewCurrencyService() *CurrencyService {
	return &CurrencyService{}
}

// GetCurrency returns the currency for the given symbol
func (cs *CurrencyService) GetCurrency(symbol string) *domain.Currency {
	return availableCurrencies[symbol]
}

// GetAvailableCurrencies returns available currencies implemented this service
func (cs *CurrencyService) GetAvailableCurrencies() []domain.Currency {
	currs := make([]domain.Currency, len(availableCurrencies))
	for _, c := range availableCurrencies {
		currs = append(currs, domain.Currency{
			Symbol:  c.Symbol,
			Decimal: c.Decimal,
		})
	}

	return currs
}
