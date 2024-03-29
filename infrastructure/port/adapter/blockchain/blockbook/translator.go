package blockbook

import (
	"fmt"
	"math/big"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchain"
)

// BitcoinTranslator is a translator for Blockbook API
type BitcoinTranslator struct{}

// ToAccountMovements converts data returning from third-party service to domain.AccountMovement value
func (tr BitcoinTranslator) ToAccountMovements(address string, v interface{}) (*domain.AccountMovements, error) {
	txs, _ := v.([]Transaction)
	am := domain.NewAccountMovements(address)

	for _, tx := range txs {
		// Inputs will be reflected as a spent
		for _, in := range tx.Inputs {
			if in.Addresses[0] != address {
				continue
			}

			val, ok := new(big.Int).SetString(in.Value, 10)
			if !ok {
				return nil, fmt.Errorf("bitcoin translation error, cannot convert in.Value(%s) to bigint", in.Value)
			}
			am.Spend(tx.BlockHeight, tx.BlockTime, tx.TxID, val, "")
		}

		// Outputs will be reflected as a receive
		for _, out := range tx.Outputs {
			if out.Addresses[0] != address {
				continue
			}

			val, ok := new(big.Int).SetString(out.Value, 10)
			if !ok {
				return nil, fmt.Errorf("bitcoin translation error, cannot convert out.Value(%s) to bigint", out.Value)
			}
			am.Receive(tx.BlockHeight, tx.BlockTime, tx.TxID, val, "")
		}
	}

	return am, nil
}

// EthereumTranslator is a translator for Blockbook API
type EthereumTranslator struct{}

// ToAccountMovements converts data returning from third-party service to domain.AccountMovement value
func (tr EthereumTranslator) ToAccountMovements(address string, v interface{}) (*domain.AccountMovements, error) {
	txs, _ := v.([]Transaction)
	am := domain.NewAccountMovements(address)
	address = blockchain.NormalizeEthereumAddress(address)

	for _, tx := range txs {
		// Do not include reverted/failed transactions
		if tx.EthereumSpecific.Status != transactionStatusSuccess {
			continue
		}

		val, ok := new(big.Int).SetString(tx.Value, 10)
		if !ok {
			return nil, fmt.Errorf("blockbook ethereum translation error, cannot convert in.Value(%s) to bigint", tx.Value)
		}

		// Any value transfers from this address will be reflected as a spent
		if blockchain.NormalizeEthereumAddress(tx.Inputs[0].Addresses[0]) == address {
			am.Spend(tx.BlockHeight, tx.BlockTime, tx.TxID, val, tx.Outputs[0].Addresses[0])
		}

		// Any value transfers to this address will be reflected as a receive
		if blockchain.NormalizeEthereumAddress(tx.Outputs[0].Addresses[0]) == address {
			am.Receive(tx.BlockHeight, tx.BlockTime, tx.TxID, val, tx.Inputs[0].Addresses[0])
		}
	}

	return am, nil
}
