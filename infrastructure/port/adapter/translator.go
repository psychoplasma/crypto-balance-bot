package adapter

import domain "github.com/psychoplasma/crypto-balance-bot"

// Translator interface declares translation from outside service output to currency service inputs
type Translator interface {
	ToAccountMovements(address string, addressData interface{}) *domain.AccountMovement
}
