package cryptobot

import (
	"math/big"
)

// BalanceChange represents a change in balance
type BalanceChange struct {
	Amount *big.Int
	TxHash string
}

// AccountMovement represents the total change
// made to Account in a certain time range
type AccountMovement struct {
	Address  string
	Currency Currency
	Changes  map[int][]*BalanceChange // All the balance changes made to Account at a certain block height, map[blockheight]balanceChanges
}

// FromAccountMovement set the given AccountMovement's currency and return it
func FromAccountMovement(c Currency, am *AccountMovement) *AccountMovement {
	am.Currency = c
	return am
}

// NewAccountMovement creates a new instance of AccountMovement
func NewAccountMovement(address string) *AccountMovement {
	return &AccountMovement{
		Address: address,
		Changes: make(map[int][]*BalanceChange),
	}
}

// AddBalanceChange adds a balance change to the list of changes at the given block height
func (am *AccountMovement) AddBalanceChange(blockHeight int, txHash string, amount *big.Int) {
	if am.Changes[blockHeight] == nil {
		am.Changes[blockHeight] = make([]*BalanceChange, 0)
	}

	am.Changes[blockHeight] = append(
		am.Changes[blockHeight],
		&BalanceChange{
			Amount: amount,
			TxHash: txHash,
		})
}

// Account in a value object (even it seems like an entity).
// Only balance and index properties change over time
// but they are actually for tracking purposes which is not
// really Address's property.
type Account struct {
	address     string
	balance     *big.Int
	blockHeight int
}

// NewAccount creates a new instance of Account with the given address
func NewAccount(address string) *Account {
	return &Account{
		address:     address,
		balance:     big.NewInt(0),
		blockHeight: -1,
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

// BlockHeight returns the last block height balance is updated
func (a *Account) BlockHeight() int {
	return a.blockHeight
}

// Apply applies a movement to the current state of this account
// Movements in a AccountMovements object should be descending-ordered
// by block height. Otherwise after the first Movement applied, the
// remaining will be ignored.
func (a *Account) Apply(am *AccountMovement) {
	if am == nil || am.Address != a.address {
		return
	}

	for blockHeight, cs := range am.Changes {
		if am == nil || blockHeight <= a.blockHeight {
			return
		}

		for _, c := range cs {
			a.balance = a.balance.Add(a.balance, c.Amount)
		}

		a.blockHeight = blockHeight
	}
}
