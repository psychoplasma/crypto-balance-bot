package blockchaindotcom

import (
	"math/big"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

// BitcoinTranslator is a translator for Blockchain.com API
type BitcoinTranslator struct{}

// ToAccountMovements converts data returning from third-party service to .AccountMovement domain object
func (bt BitcoinTranslator) ToAccountMovements(address string, v interface{}) *domain.AccountMovements {
	txs, _ := v.([]Transaction)
	am := domain.NewAccountMovements(address)

	for _, tx := range txs {
		// Inputs will be reflected as a decrease in balance
		for _, in := range tx.Inputs {
			if in.PrevOutput.Address != address {
				continue
			}

			am.AddBalanceChange(tx.BlockHeight, tx.Hash, new(big.Int).Neg(in.PrevOutput.Value))
		}

		// Outputs will be reflected as an increase in balance
		for _, out := range tx.Outputs {
			if out.Address != address {
				continue
			}

			am.AddBalanceChange(tx.BlockHeight, tx.Hash, out.Value)
		}
	}

	return am
}
