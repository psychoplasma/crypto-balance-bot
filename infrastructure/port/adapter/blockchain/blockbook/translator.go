package blockbook

import (
	"fmt"
	"math/big"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

// BitcoinTranslator is a translator for Blockbook API
type BitcoinTranslator struct{}

// ToAccountMovements converts data returning from third-party service to .AccountMovement domain object
func (bt BitcoinTranslator) ToAccountMovements(address string, v interface{}) (*domain.AccountMovements, error) {
	txs, _ := v.([]Transaction)
	am := domain.NewAccountMovements(address)

	for _, tx := range txs {
		// Inputs will be reflected as a decrease in balance
		for _, in := range tx.Inputs {
			if in.Addresses[0] != address {
				continue
			}

			val, ok := new(big.Int).SetString(in.Value, 10)
			if !ok {
				return nil, fmt.Errorf("bitcoin translation error, cannot convert in.Value(%s) to bigint", in.Value)
			}
			am.AddBalanceChange(tx.BlockHeight, tx.TxID, new(big.Int).Neg(val))
		}

		// Outputs will be reflected as an increase in balance
		for _, out := range tx.Outputs {
			if out.Addresses[0] != address {
				continue
			}

			val, ok := new(big.Int).SetString(out.Value, 10)
			if !ok {
				return nil, fmt.Errorf("bitcoin translation error, cannot convert out.Value(%s) to bigint", out.Value)
			}
			am.AddBalanceChange(tx.BlockHeight, tx.TxID, val)
		}
	}

	return am, nil
}
