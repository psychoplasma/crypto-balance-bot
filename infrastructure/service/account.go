package service

import (
	"log"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

type AccountService struct {
	cs domain.CurrencyService
}

func NewAccountService(cs domain.CurrencyService) *AccountService {
	return &AccountService{
		cs: cs,
	}
}

// FetchAccountMovements fetches balances in this account and returns any balance changes
func (as *AccountService) FetchAccountMovements(a *domain.Account) ([]*domain.AccountMovement, error) {
	m, err := as.cs.GetAddressTxs(a.Address(), a.BlockHeight())
	if err != nil {
		log.Printf("failed to fetch movements for address(%s), %s", a.Address(), err)
		return nil, err
	}

	return m, nil
}
