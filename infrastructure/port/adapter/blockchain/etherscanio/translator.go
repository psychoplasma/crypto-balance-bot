package etherscanio

import (
	"fmt"
	"math/big"
	"strconv"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchain"
)

// EthereumTranslator is a translator for Etherscan.io API
type EthereumTranslator struct{}

// ToAccountMovements converts data returning from third-party service to AccountMovement domain object
func (et EthereumTranslator) ToAccountMovements(address string, v interface{}) (*domain.AccountMovements, error) {
	txs, _ := v.([]Transaction)
	am := domain.NewAccountMovements(address)
	address = blockchain.NormalizeEthereumAddress(address)

	for _, tx := range txs {
		// Do not include reverted/failed transactions
		if tx.Status != transactionStatusSuccess {
			continue
		}

		blockHeight, err := strconv.ParseInt(tx.BlockHeight, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("etherscanio ethereum transalation error, %s", err.Error())
		}

		val, ok := new(big.Int).SetString(tx.Value, 10)
		if !ok {
			return nil, fmt.Errorf("etherscanio ethereum transalation error, cannot convert tx.Value(%s) to bigint", tx.Value)
		}

		// Any value transfers from this address will be reflected as a decrease in balance
		if blockchain.NormalizeEthereumAddress(tx.From) == address {
			am.AddBalanceChange(int(blockHeight), tx.Hash, new(big.Int).Neg(val))
		}

		// Any value transfers to this address will be reflected as an increase in balance
		if blockchain.NormalizeEthereumAddress(tx.To) == address {
			am.AddBalanceChange(int(blockHeight), tx.Hash, val)
		}
	}

	return am, nil
}
