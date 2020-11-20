package main

import (
	"github.com/psychoplasma/crypto-balance-bot/application"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/api/coin"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/persistance/inmemory"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/service"
)

var subsRepo = inmemory.NewSubscriptionReposititory()
var subsAppService = application.NewSubscriptionApplication(subsRepo)
var currencyAppService = application.NewCurrencyService()
var obs = application.NewObserver(service.NewAccountService(&coin.BitcoinAPI{}), subsRepo)

func main() {
	b := NewBot(subsAppService, currencyAppService, obs)
	b.RegisterCommands()

	go obs.Observe()

	b.Start()
}
