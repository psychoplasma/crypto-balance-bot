package cryptobot

import (
	"math/big"
	"sort"
	"time"
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

// AccountAssetsMovedEvent represents a domain event upon AccountMovements
type AccountAssetsMovedEvent struct {
	version    int
	occurredOn time.Time
	subsID     string
	acms       *AccountMovements
	c          Currency
}

// NewAccountAssetsMovedEvent creates a new instance from AccountMovements
func NewAccountAssetsMovedEvent(subsID string, c Currency, acms *AccountMovements) *AccountAssetsMovedEvent {
	return &AccountAssetsMovedEvent{
		version:    1,
		occurredOn: time.Now(),
		subsID:     subsID,
		acms:       acms,
		c:          c,
	}
}

// AccountMovements returns AccountMovements
func (evt *AccountAssetsMovedEvent) AccountMovements() *AccountMovements {
	return evt.acms
}

// Currency returns the currency property
func (evt *AccountAssetsMovedEvent) Currency() Currency {
	return evt.c
}

// SubscriptionID returns the subsID property
func (evt *AccountAssetsMovedEvent) SubscriptionID() string {
	return evt.subsID
}

// OccurredOn returns event time
func (evt *AccountAssetsMovedEvent) OccurredOn() time.Time {
	return evt.occurredOn
}

// EventVersion returns event version
func (evt *AccountAssetsMovedEvent) EventVersion() int {
	return evt.version
}
