package main

import (
	"github.com/psychoplasma/crypto-balance-bot/application"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/persistance/inmemory"
)

var subsRepo = inmemory.NewSubscriptionReposititory()
var subsAppService = application.NewSubscriptionApplication(subsRepo)
var currencyAppService = application.NewCurrencyService()

func main() {
	b := NewBot(subsAppService, currencyAppService)
	b.RegisterCommands()
	b.Start()
}
