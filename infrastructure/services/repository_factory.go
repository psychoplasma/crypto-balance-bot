package services

import (
	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/persistence/mongodb"
)

// RepositoryServiceFactory keeps implemented repository services
var RepositoryServiceFactory = map[string]domain.SubscriptionRepository{
	"mongodb":    mongodb.NewSubscriptionRepository(),
	"postgresql": nil,
}
