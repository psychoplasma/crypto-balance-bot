package cryptobot

import (
	"math/big"
	"sort"
	"time"
)

// Type of balance change
const (
	ReceivedBalance = iota
	SpentBalance
)

// BalanceChange represents a change in balance
type balanceChange struct {
	Type   int
	Amount *big.Int
	TxHash string
}

func (bc balanceChange) Value() *big.Int {
	switch bc.Type {
	case ReceivedBalance:
		return big.NewInt(bc.Amount.Int64())
	case SpentBalance:
		return new(big.Int).Neg(bc.Amount)
	default:
		return big.NewInt(0)
	}
}

// AccountMovements represents the total change
// made to Account in a certain time range
type AccountMovements struct {
	Address string
	Blocks  []uint64
	Changes map[uint64][]*balanceChange
}

// NewAccountMovements creates a new instance of AccountMovement
func NewAccountMovements(address string) *AccountMovements {
	return &AccountMovements{
		Address: address,
		Changes: make(map[uint64][]*balanceChange, 0),
		Blocks:  []uint64{},
	}
}

// Sort sorts the changes by block height in ascending order
func (am *AccountMovements) Sort() *AccountMovements {
	// Sort am.Blocks in increasing order
	sort.Slice(am.Blocks, func(i, j int) bool { return am.Blocks[i] < am.Blocks[j] })
	return am
}

// ReceiveBalance adds a balance change as received to the list of changes at the given block height
func (am *AccountMovements) ReceiveBalance(blockHeight uint64, txHash string, amount *big.Int) {
	am.addBalanceChange(blockHeight, txHash, amount, ReceivedBalance)
}

// SpendBalance adds a balance change as spent to the list of changes at the given block height
func (am *AccountMovements) SpendBalance(blockHeight uint64, txHash string, amount *big.Int) {
	am.addBalanceChange(blockHeight, txHash, amount, SpentBalance)
}

func (am *AccountMovements) addBalanceChange(blockHeight uint64, txHash string, amount *big.Int, bType int) {
	if am.Changes[blockHeight] == nil {
		am.Blocks = append(am.Blocks, blockHeight)
		am.Changes[blockHeight] = make([]*balanceChange, 0)
	}

	am.Changes[blockHeight] = append(
		am.Changes[blockHeight],
		&balanceChange{
			Amount: amount,
			TxHash: txHash,
			Type:   bType,
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
