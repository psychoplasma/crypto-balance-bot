package blockchaindotcom

import (
	domain "github.com/psychoplasma/crypto-balance-bot"
)

// BitcoinTranslator is a translator for Blockchain.com API
type BitcoinTranslator struct{}

// ToAccountMovements converts data returning from third-party service to .AccountMovement domain object
func (bt BitcoinTranslator) ToAccountMovements(address string, v interface{}) (*domain.AccountMovements, error) {
	txs, _ := v.([]Transaction)
	am := domain.NewAccountMovements(address)

	for _, tx := range txs {
		// Inputs will be reflected as a spent
		for _, in := range tx.Inputs {
			if in.PrevOutput.Address != address {
				continue
			}

			am.SpendBalance(tx.BlockHeight, tx.Hash, in.PrevOutput.Value)
		}

		// Outputs will be reflected as a receive
		for _, out := range tx.Outputs {
			if out.Address != address {
				continue
			}

			am.ReceiveBalance(tx.BlockHeight, tx.Hash, out.Value)
		}
	}

	return am, nil
}
