package coin

import (
	"math/big"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

// EthereumAPI implements CurrencyAPI for Ethereum
type EthereumAPI struct {
}

// GetBalance fetches balance of the given address
func (a *EthereumAPI) GetBalance(address string) (*big.Int, error) {
	return nil, nil
}

// GetTransactions fetches transaction of the given address starting from the given index
func (a *EthereumAPI) GetTransactions(address string, since int) ([]*domain.Transaction, error) {
	return nil, nil
}

// ValidateAddress checks whether or not the given address is valid
// and returns an error in case of invalid address
func (a *EthereumAPI) ValidateAddress(address string) error {
	return nil
}
