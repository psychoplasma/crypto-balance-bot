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

		blockHeight, err := strconv.ParseUint(tx.BlockHeight, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("etherscanio ethereum transalation error, %s", err.Error())
		}

		val, ok := new(big.Int).SetString(tx.Value, 10)
		if !ok {
			return nil, fmt.Errorf("etherscanio ethereum transalation error, cannot convert tx.Value(%s) to bigint", tx.Value)
		}

		from := blockchain.NormalizeEthereumAddress(tx.From)
		to := blockchain.NormalizeEthereumAddress(tx.To)
		timestamp, _ := strconv.ParseUint(tx.Timestamp, 10, 64)

		// Any value transfers from this address will be reflected as a spent
		if from == address {
			am.Spend(blockHeight, timestamp, tx.Hash, val, to)
		}

		// Any value transfers to this address will be reflected as a receive
		if to == address {
			am.Receive(blockHeight, timestamp, tx.Hash, val, from)
		}
	}

	return am, nil
}
