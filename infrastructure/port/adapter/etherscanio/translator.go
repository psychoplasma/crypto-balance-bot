package etherscanio

import (
	"math/big"
	"strconv"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

// EthereumTranslator is a translator for Etherscan.io API
type EthereumTranslator struct{}

// ToAccountMovements converts data returning from third-party service to AccountMovement domain object
func (et EthereumTranslator) ToAccountMovements(address string, v interface{}) []*domain.AccountMovement {
	txs, _ := v.([]Transaction)
	ams := []*domain.AccountMovement{}

	for _, tx := range txs {
		// Do not include reverted/failed transactions
		if tx.Status == "1" {
			continue
		}

		am := &domain.AccountMovement{}
		i, _ := strconv.ParseInt(tx.BlockHeight, 10, 64)

		am.BlockHeight = int(i)
		// Any value transfers from this address will be reflected as a decrease in balance
		if tx.From == address {
			a, _ := big.NewInt(0).SetString(tx.Value, 10)
			ch := &domain.BalanceChange{
				Amount: big.NewInt(0).Neg(a),
				TxHash: tx.Hash,
			}

			am.Changes = append(am.Changes, ch)
		}

		// Any value transfers to this address will be reflected as an increase in balance
		if tx.To == address {
			a, _ := big.NewInt(0).SetString(tx.Value, 10)
			ch := &domain.BalanceChange{
				Amount: a,
				TxHash: tx.Hash,
			}

			am.Changes = append(am.Changes, ch)
		}

		ams = append(ams, am)
	}

	return ams
}
