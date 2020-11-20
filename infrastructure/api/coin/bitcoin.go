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

// GetAddressTxs fetches txs of the given address since the given block height(exclusive)
func (a *BitcoinAPI) GetAddressTxs(address string, sinceBlockHeight int) ([]*domain.AccountMovement, error) {
	test := []*domain.AccountMovement{
		{
			BlockHeight: 1,
			Changes: []*domain.BalanceChange{
				{
					Amount: big.NewInt(5),
				},
			},
		},
		{
			BlockHeight: 2,
			Changes: []*domain.BalanceChange{
				{
					Amount: big.NewInt(-3),
				},
			},
		},
		{
			BlockHeight: 3,
			Changes: []*domain.BalanceChange{
				{
					Amount: big.NewInt(8),
				},
			},
		},
	}
	return test, nil
}

// ValidateAddress checks whether or not the given address is valid
// and returns an error in case of invalid address
func (a *BitcoinAPI) ValidateAddress(address string) error {
	return nil
}
