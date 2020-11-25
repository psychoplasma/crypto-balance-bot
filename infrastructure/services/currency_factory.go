package services

import (
	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchaindotcom"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/etherscanio"
)

// Implemented currencies
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

// CurrencyFactory keeps implemented currencies
var CurrencyFactory = map[string]*domain.Currency{
	"btc": &BTC,
	"eth": &ETH,
}

// CurrencyServiceFactory keeps implemented currency services
var CurrencyServiceFactory = map[string]domain.CurrencyService{
	"btc": blockchaindotcom.NewBitcoinAPI(blockchaindotcom.BitcoinTranslator{}),
	"eth": etherscanio.NewEthereumAPI(etherscanio.EthereumTranslator{}),
}

// Translator interface declares translation from outside service output to currency service inputs
type Translator interface {
	ToAccountMovements(address string, addressData interface{}) []*domain.AccountMovement
}
