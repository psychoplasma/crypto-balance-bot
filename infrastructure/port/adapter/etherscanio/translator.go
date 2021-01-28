package etherscanio

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

const addressPrefix = "0x"

// EthereumTranslator is a translator for Etherscan.io API
type EthereumTranslator struct{}

// ToAccountMovements converts data returning from third-party service to AccountMovement domain object
func (et EthereumTranslator) ToAccountMovements(address string, v interface{}) (*domain.AccountMovements, error) {
	txs, _ := v.([]Transaction)
	am := domain.NewAccountMovements(address)

	for _, tx := range txs {
		// Do not include reverted/failed transactions
		if tx.Status != transactionStatusSuccess {
			continue
		}

		blockHeight, err := strconv.ParseInt(tx.BlockHeight, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("ethereum transalation error, %s", err.Error())
		}

		// Any value transfers from this address will be reflected as a decrease in balance
		if normalizeAddress(tx.From) == normalizeAddress(address) {
			a, ok := new(big.Int).SetString(tx.Value, 10)
			if !ok {
				return nil, fmt.Errorf("ethereum transalation error, cannot convert tx.Value to bigint")
			}
			am.AddBalanceChange(int(blockHeight), tx.Hash, new(big.Int).Neg(a))
		}

		// Any value transfers to this address will be reflected as an increase in balance
		if normalizeAddress(tx.To) == normalizeAddress(address) {
			a, ok := new(big.Int).SetString(tx.Value, 10)
			if !ok {
				return nil, fmt.Errorf("ethereum transalation error, cannot convert tx.Value to bigint")
			}
			am.AddBalanceChange(int(blockHeight), tx.Hash, a)
		}
	}

	return am, nil
}

func normalizeAddress(address string) string {
	address = strings.Trim(address, " ")

	if !strings.HasPrefix(address, addressPrefix) {
		address = addressPrefix + address
	}

	return strings.ToLower(address)
}
