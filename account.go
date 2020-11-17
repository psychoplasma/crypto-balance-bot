package cryptobot

import (
	"errors"
	"log"
	"math/big"
)

// Represets error related to account operations
var (
	ErrInvalidAddrDesc = errors.New("invalid address descriptor")
)

// Account represents Account to be subscribed to bot
type Account struct {
	c           Currency
	xPubKey     string
	AddressList []string
	Balances    map[string]*big.Int // Keeps balances for each address
	TxHashes    map[string][]string // Keeps transaction hashes for each address
}

// NewAccount creates a new instance of Account with the given address descriptor and currency
func NewAccount(c Currency, addrDesc string) (*Account, error) {
	if err := c.API.ValidatePubKey(addrDesc); err == nil {
		return newAccountByMasterPubKey(c, addrDesc), nil
	}

	if err := c.API.ValidateAddress(addrDesc); err == nil {
		return newAccountByAddress(c, addrDesc), nil
	}

	return nil, ErrInvalidAddrDesc
}

// UpdateBalances updates balances in this account and returns any balance change
func (a *Account) UpdateBalances() map[string]*big.Int {
	movements := map[string]*big.Int{}
	for _, addr := range a.AddressList {
		b, err := a.c.API.GetBalance(addr)
		if err != nil {
			log.Printf("cannot fetch balance for address(%s), %s", addr, err)
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

// UpdateTxs updates tx hashes for each address
// in this account and returns the changes
func (a *Account) UpdateTxs() map[string][]string {
	changes := map[string][]string{}
	for _, addr := range a.AddressList {
		txs, err := a.c.API.GetTransactions(addr, len(a.TxHashes[addr]))
		if err != nil {
			log.Printf("cannot fetch transactions for address(%s), %s", addr, err)
		}

		if len(txs) > 0 {
			changes[addr] = txs
		}
		a.TxHashes[addr] = append(a.TxHashes[addr], txs...)
	}

	return changes
}

// NewAccountByAddress creates a new instance of Account with the given address
func newAccountByAddress(c Currency, address string) *Account {
	a := &Account{
		c:           c,
		AddressList: []string{address},
		Balances:    map[string]*big.Int{},
		TxHashes:    map[string][]string{},
	}
	return a
}

// NewAccountByMasterPubKey creates a new instance of Account
// which consist of addresses drived from the given master public key
func newAccountByMasterPubKey(c Currency, xPubKey string) *Account {
	a := &Account{
		c:           c,
		xPubKey:     xPubKey,
		AddressList: deriveAddresses(xPubKey),
		Balances:    map[string]*big.Int{},
		TxHashes:    map[string][]string{},
	}
	return a
}

// TODO: implement
func deriveAddresses(xPubKey string) []string {
	return []string{}
}
