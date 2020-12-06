package cryptobot

import (
	"math/big"
	"sort"
)

// BalanceChange represents a change in balance
type balanceChange struct {
	Amount *big.Int
	TxHash string
}

// AccountMovements represents the total change
// made to Account in a certain time range
type AccountMovements struct {
	Address string
	Blocks  []int
	Changes map[int][]*balanceChange
}

// NewAccountMovements creates a new instance of AccountMovement
func NewAccountMovements(address string) *AccountMovements {
	return &AccountMovements{
		Address: address,
		Changes: make(map[int][]*balanceChange, 0),
		Blocks:  []int{},
	}
}

// AccountMovementsFrom is a acopy constructor
func AccountMovementsFrom(acm *AccountMovements) *AccountMovements {
	return &AccountMovements{
		Address: acm.Address,
		Changes: acm.Changes,
		Blocks:  acm.Blocks,
	}
}

// Sort sorts the changes by block height in ascending order
func (am *AccountMovements) Sort() *AccountMovements {
	sort.Ints(am.Blocks)
	return am
}

// AddBalanceChange adds a balance change to the list of changes at the given block height
func (am *AccountMovements) AddBalanceChange(blockHeight int, txHash string, amount *big.Int) {
	if am.Changes[blockHeight] == nil {
		am.Blocks = append(am.Blocks, blockHeight)
		am.Changes[blockHeight] = make([]*balanceChange, 0)
	}

	am.Changes[blockHeight] = append(
		am.Changes[blockHeight],
		&balanceChange{
			Amount: amount,
			TxHash: txHash,
		},
	)
}

// SubscriptionMovements keeps track of every asset movements under each account of the subscription
// Blockheight ranges for movements for different accounts may differ
type SubscriptionMovements struct {
	subsID string
	c      Currency
	acms   map[string]*AccountMovements
}

// NewSubscriptionMovements creates a new instance from the given subscription ID and Currency
func NewSubscriptionMovements(id string, c Currency) *SubscriptionMovements {
	return &SubscriptionMovements{
		subsID: id,
		c:      c,
		acms:   make(map[string]*AccountMovements),
	}
}

// Currency returns the Currency property
func (sm *SubscriptionMovements) Currency() Currency {
	return sm.c
}

// Subscription returns the Subscription ID property
func (sm *SubscriptionMovements) Subscription() string {
	return sm.subsID
}

// AccountMovements retunrs all the movements of each account for this subscription
func (sm *SubscriptionMovements) AccountMovements() map[string]*AccountMovements {
	for _, acm := range sm.acms {
		acm.Sort()
	}
	return sm.acms
}

// AccountMovementsForAccount retunrs the movement for the given address of an account
func (sm *SubscriptionMovements) AccountMovementsForAccount(address string) *AccountMovements {
	return sm.acms[address].Sort()
}

// AddAccountMovements adds a set of account movements
// If the given account movements alread exist for
// the given account ignores the movements
func (sm *SubscriptionMovements) AddAccountMovements(acm *AccountMovements) {
	if len(acm.Changes) == 0 {
		return
	}

	if _, exist := sm.acms[acm.Address]; exist {
		return
	}

	sm.acms[acm.Address] = AccountMovementsFrom(acm)
}

// Account in a value object (even it seems like an entity).
// Only balance and index properties change over time
// but they are actually for tracking purposes which is not
// really Address's property.
type Account struct {
	address     string
	balance     *big.Int
	currency    Currency
	blockHeight int
}

// NewAccount creates a new instance of Account with the given address
func NewAccount(address string, c Currency) *Account {
	return &Account{
		address:     address,
		balance:     big.NewInt(0),
		currency:    c,
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

// Currency returns the currency of this account
func (a *Account) Currency() Currency {
	return a.currency
}

// Apply applies a movement to the current state of this account
// Movements in a AccountMovements object should be descending-ordered
// by block height. Otherwise after the first Movement applied, the
// remaining will be ignored.
func (a *Account) Apply(am *AccountMovements) {
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
