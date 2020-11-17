package main

import (
	"github.com/psychoplasma/crypto-balance-bot/cmd/application"
	"github.com/psychoplasma/crypto-balance-bot/repo"
)

var subsRepo = repo.SubscriptionInMemoryReposititory()
var subsAppService = application.SubscriptionService(subsRepo)
var currencyAppService = application.CurrencyService()

func main() {

}
