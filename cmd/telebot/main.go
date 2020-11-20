package main

import (
	"github.com/psychoplasma/crypto-balance-bot/application"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/persistance/inmemory"
)

var subsAppService = application.NewSubscriptionApplication(
	inmemory.NewSubscriptionReposititory())
var currencyAppService = application.NewCurrencyService()

func main() {
	b := NewBot(subsAppService, currencyAppService)
	b.RegisterCommands()
	b.Start()
}
