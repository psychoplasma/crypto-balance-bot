package main

import (
	"github.com/psychoplasma/crypto-balance-bot/application"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/notification"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/persistence/inmemory"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter"
)

var subsRepo = inmemory.NewSubscriptionReposititory()
var movementPublisher = adapter.NewTelegramPublisher(
	"873977886:AAEJetV4LiotkaqDo3NGOrZXQ2BWEA2U8ts", notification.MovementFormatter{})

func main() {
	o := application.NewObserver(subsRepo)
	o.RegisterPublisher(movementPublisher)
	o.Observe()
}
