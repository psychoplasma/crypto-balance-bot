package cryptobot

import (
	"log"
	"math/big"
)

// CurrencyAPI represents API to fetch relavent info about account for the given currency
type CurrencyAPI interface {
	GetBalance(addressDesc string) (*big.Int, error)
}

// Account represents Account to be subscribed to bot
type Account struct {
	CurrencyID   string
	MasterPubKey string
	AddressList  []string
	Balances     map[string]*big.Int
	API          CurrencyAPI
}

// NewAccountByAddress creates a new instance of Account with the given address
func NewAccountByAddress(currencyID string, address string, api CurrencyAPI) *Account {
	a := &Account{
		CurrencyID:  currencyID,
		AddressList: []string{address},
		Balances:    map[string]*big.Int{},
		API:         api,
	}
	a.UpdateBalances()
	return a
}

// NewAccountByMasterPubKey creates a new instance of Account which consist of addresses drived from the given master public key
func NewAccountByMasterPubKey(currencyID string, masterPubKey string, api CurrencyAPI) *Account {
	a := &Account{
		CurrencyID:   currencyID,
		MasterPubKey: masterPubKey,
		AddressList:  deriveAddresses(masterPubKey),
		Balances:     map[string]*big.Int{},
		API:          api,
	}
	a.UpdateBalances()
	return a
}

// UpdateBalances updates balances in this account and returns any balance change
func (a *Account) UpdateBalances() map[string]*big.Int {
	movements := map[string]*big.Int{}
	for _, addr := range a.AddressList {
		b, err := a.API.GetBalance(addr)
		if err != nil {
			log.Printf("cannot fetch balance, %s", err)
		}

		diff := big.NewInt(0)
		diff = diff.Sub(b, a.Balances[addr])

		if diff.Cmp(big.NewInt(0)) != 0 {
			movements[addr] = diff
		}
		a.Balances[addr] = b
	}

	return movements
}

// TODO: implement
func deriveAddresses(masterPubKey string) []string {
	return []string{}
}
