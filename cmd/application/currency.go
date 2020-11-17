package application

import (
	cryptoBot "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/coins"
)

// AvailableCurrencies contains the implemented currencies
// TODO: Load these through a configuration file
var availableCurrencies = map[string]*cryptoBot.Currency{
	"btc": {
		Decimal: 8,
		Symbol:  "btc",
		API:     &coins.BitcoinAPI{},
	},

	"eth": {
		Decimal: 16,
		Symbol:  "eth",
		API:     &coins.EthereumAPI{},
	},
}

// Currency represents presentation model for Currency entity in our domain
type Currency struct {
	Symbol  string `json:"symbol"`
	Decimal int    `json:"decimal"`
}

// CurrencyService exposes application services for currency entity
type CurrencyService struct {
}

// CurrencyService factory function
func CurrencyService() *CurrencyService {
	return &CurrencyService{}
}

// GetCurrency returns the currency for the given symbol
func (cs *CurrencyService) GetCurrency(symbol string) Currency {
	return availableCurrencies[symbol]
}

// GetAvailableCurrencies returns available currencies implemented this service
func (cs *CurrencyService) GetAvailableCurrencies() []Currency {
	currs := make([]Currency, len(availableCurrencies))
	for _, c := range availableCurrencies {
		currs = append(currs, Currency{
			Symbol:  c.Symbol,
			Decimal: c.Decimal,
		})
	}

	return currs
}
