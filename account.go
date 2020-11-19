package cryptobot

import (
	"errors"
	"fmt"
	"math/big"
)

// Represets error related to account operations
var (
	ErrInvalidAddrDesc = errors.New("invalid address descriptor")
	ErrInvalidXPubKey  = errors.New("unable to derive addresses from the given master public key")
)

type AccountService interface {
	UpdateBalances(a *Account) (map[string]*big.Int, error)
	UpdateTxs(a *Account) (map[string][]string, error)
}

// Transaction is a value object
type Transaction struct {
	Hash         string
	BlockHeight  int
	changeAmount *big.Int
}

type Balance struct {
	Amount *big.Int
	c      Currency
}

func (b Balance) ToString() string {
	return fmt.Sprintf("%s %s", b.Amount.Text(10), b.c.Symbol)
}

// Account in a value object (even it seems like an entity).
// Only balance and tx related properties change over time
// but they are actually for tracking purposes which is not
// really a Address's property.
type Account struct {
	address string
	b       *Balance
	txCount int
	txs     map[string]*Transaction
}

// NewAccount creates a new instance of Account with the given address
func NewAccount(c Currency, address string) *Account {
	return &Account{
		address: address,
		b: &Balance{
			Amount: big.NewInt(0),
			c:      c,
		},
	}
}

func (a *Account) Address() string {
	return a.address
}

func (a *Account) Balance() Balance {
	return *a.b
}

func (a *Account) TxCount() int {
	return a.txCount
}

func (a *Account) AddTxs(txs []*Transaction) {
	for _, tx := range txs {
		if a.txs[tx.Hash] != nil {
			continue
		}

		a.txs[tx.Hash] = tx
		a.txCount++
	}
}

func (a *Account) UpdateBalance(b *big.Int) {
	a.b.Amount = b
}
