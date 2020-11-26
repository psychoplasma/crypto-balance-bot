package cryptobot

import (
	"fmt"
	"math/big"
)

// AccountMovement represents the total change
// made to Account in a certain block height
type AccountMovement struct {
	BlockHeight int
	Changes     []*BalanceChange
}

// BalanceChange represents a change in balance
type BalanceChange struct {
	Amount *big.Int
	TxHash string
}

// Account in a value object (even it seems like an entity).
// Only balance and index properties change over time
// but they are actually for tracking purposes which is not
// really Address's property.
type Account struct {
	address     string
	balance     *big.Int
	blockHeight int
	c           Currency
}

// NewAccount creates a new instance of Account with the given address
func NewAccount(c Currency, address string) *Account {
	return &Account{
		address:     address,
		balance:     big.NewInt(0),
		blockHeight: -1,
		c:           c,
	}
}

// Address returns address of this account
func (a *Account) Address() string {
	return a.address
}

// Balance returns the last checked balance
func (a *Account) Balance() *big.Int {
	return a.balance
}

// BalanceToString returns the string representation of balance
func (a *Account) BalanceToString() string {
	return fmt.Sprintf("%s %s", a.balance.Text(10), a.c.Symbol)
}

// BlockHeight returns the last block height balance is updated
func (a *Account) BlockHeight() int {
	return a.blockHeight
}

// Currency returns the currency type of this account
func (a *Account) Currency() Currency {
	return a.c
}

// Apply applies a movement to the current state of this account
func (a *Account) Apply(am *AccountMovement) {
	if am == nil || am.BlockHeight <= a.blockHeight {
		return
	}

	for _, c := range am.Changes {
		a.balance = a.balance.Add(a.balance, c.Amount)
	}

	a.blockHeight = am.BlockHeight
}
