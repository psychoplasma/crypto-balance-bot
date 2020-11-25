package blockchaindotcom

import (
	"math/big"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

// BitcoinTranslator is a translator for Blockchain.com API
type BitcoinTranslator struct{}

// ToAccountMovements converts data returning from third-party service to .AccountMovement domain object
func (bt BitcoinTranslator) ToAccountMovements(address string, v interface{}) []*domain.AccountMovement {
	txs, _ := v.([]Transaction)
	ams := []*domain.AccountMovement{}

	for _, tx := range txs {
		am := &domain.AccountMovement{}
		am.BlockHeight = tx.BlockHeight

		// Inputs will be reflected as a decrease in balance
		for _, in := range tx.Inputs {
			if in.PrevOutput.Address != address {
				continue
			}

			ch := &domain.BalanceChange{
				Amount: big.NewInt(0).Neg(in.PrevOutput.Value),
				TxHash: tx.Hash,
			}

			am.Changes = append(am.Changes, ch)
		}

		// Outputs will be reflected as an increase in balance
		for _, out := range tx.Outputs {
			if out.Address != address {
				continue
			}

			ch := &domain.BalanceChange{
				Amount: out.Value,
				TxHash: tx.Hash,
			}

			am.Changes = append(am.Changes, ch)
		}

		if am.Changes == nil || len(am.Changes) < 1 {
			continue
		}

		ams = append(ams, am)
	}

	return ams
}
