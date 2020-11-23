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

	// Assumption made here is that there won't be more than one transaction
	// in the same block fot the given address. Another saying, each tx in Txs
	// will have different block heights. This doesn't necessarily have to be
	// the case, because there can be utxo spents from the same address in
	// the same block. However due to the fact that the purpose of this application
	// to let know the subscribers about whether or not there is a utxo spent in
	// for the given address, even if there are more than one utxo spent from
	// one address, catching only one of them will suffice.
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
