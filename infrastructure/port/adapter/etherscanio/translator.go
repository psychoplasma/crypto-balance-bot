package etherscanio

import (
	"math/big"
	"strconv"
	"strings"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

const addressPrefix = "0x"

// EthereumTranslator is a translator for Etherscan.io API
type EthereumTranslator struct{}

// ToAccountMovements converts data returning from third-party service to AccountMovement domain object
func (et EthereumTranslator) ToAccountMovements(address string, v interface{}) *domain.AccountMovement {
	txs, _ := v.([]Transaction)
	am := domain.NewAccountMovement(address)

	for _, tx := range txs {
		// Do not include reverted/failed transactions
		if tx.Status != transactionStatusSuccess {
			continue
		}

		blockHeight, _ := strconv.ParseInt(tx.BlockHeight, 10, 64)

		// Any value transfers from this address will be reflected as a decrease in balance
		if normalizeAddress(tx.From) == normalizeAddress(address) {
			a, _ := new(big.Int).SetString(tx.Value, 10)
			am.AddBalanceChange(int(blockHeight), tx.Hash, new(big.Int).Neg(a))
		}

		// Any value transfers to this address will be reflected as an increase in balance
		if normalizeAddress(tx.To) == normalizeAddress(address) {
			a, _ := new(big.Int).SetString(tx.Value, 10)
			am.AddBalanceChange(int(blockHeight), tx.Hash, a)
		}
	}

	return am
}

func normalizeAddress(address string) string {
	address = strings.Trim(address, " ")

	if !strings.HasPrefix(address, addressPrefix) {
		address = addressPrefix + address
	}

	return strings.ToLower(address)
}
