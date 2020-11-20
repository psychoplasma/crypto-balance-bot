package coin

import (
	"math/big"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

// BitcoinAPI implements CurrencyAPI for Bitcoin
type BitcoinAPI struct {
}

// GetBalance fetches balance of the given address
func (a *BitcoinAPI) GetBalance(address string) (*big.Int, error) {
	return nil, nil
}

// GetTransactions fetches transaction of the given address starting from the given index
func (a *BitcoinAPI) GetTransactions(address string, index int) ([]*domain.Transaction, error) {
	return nil, nil
}

// ValidateAddress checks whether or not the given address is valid
// and returns an error in case of invalid address
func (a *BitcoinAPI) ValidateAddress(address string) error {
	return nil
}
