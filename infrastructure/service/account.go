package service

import (
	"log"
	"math/big"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

type AccountService struct {
	cs domain.CurrencyService
}

// UpdateBalances updates balances in this account and returns any balance change
func (as *AccountService) UpdateBalances(s *domain.Subscription) (map[string]*big.Int, error) {
	changes := map[string]*big.Int{}
	for _, a := range s.Accounts {
		b, err := as.cs.GetBalance(a.Address())
		if err != nil {
			log.Printf("failed to fetch balance for address(%s), %s", a.Address(), err)
			return nil, err
		}

		diff := big.NewInt(0)
		diff = diff.Sub(b, a.Balance().Amount)

		if diff.Cmp(big.NewInt(0)) != 0 {
			changes[a.Address()] = diff
		}

		a.UpdateBalance(b)
	}

	return changes, nil
}

// UpdateTxs updates tx hashes for each address
// in this account and returns the changes
func (as *AccountService) UpdateTxs(a *domain.Subscription) (map[string][]*domain.Transaction, error) {
	changes := map[string][]*domain.Transaction{}
	for _, a := range a.Accounts {
		txs, err := as.cs.GetTransactions(a.Address(), a.TxCount())
		if err != nil {
			log.Printf("failed to fetch transactions for address(%s), %s", a.Address(), err)
			return nil, err
		}

		if len(txs) > 0 {
			changes[a.Address()] = txs
		}

		a.AddTxs(txs)
	}

	return changes, nil
}
