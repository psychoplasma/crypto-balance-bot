package blockchain

import domain "github.com/psychoplasma/crypto-balance-bot"

// Translator interface declares translation from third-party blockchain service outputs to domain service inputs
type Translator interface {
	ToAccountMovements(address string, addressData interface{}) (*domain.AccountMovements, error)
}
