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

// GetAddressTxs fetches txs of the given address since the given block height(exclusive)
func (a *EthereumAPI) GetAddressTxs(address string, sinceBlockHeight int) ([]*domain.AccountMovement, error) {
	return nil, nil
}

// ValidateAddress checks whether or not the given address is valid
// and returns an error in case of invalid address
func (a *EthereumAPI) ValidateAddress(address string) error {
	return nil
}
