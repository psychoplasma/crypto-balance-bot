package services

import (
	"math/big"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchaindotcom"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/etherscanio"
)

// Implemented currencies
var (
	BTC = domain.Currency{
		Decimal: big.NewInt(100000000),
		Symbol:  "btc",
	}
	ETH = domain.Currency{
		Decimal: big.NewInt(1000000000000000000),
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
