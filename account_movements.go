package cryptobot

import (
	"math/big"
	"sort"
	"time"
)

// Type of balance change
const (
	Received = iota
	Spent
)

// Transfer represents a change in balance
type Transfer struct {
	Type        int
	Address     string
	Amount      *big.Int
	BlockHeight uint64
	Timestamp   uint64
	TxHash      string
}

// Value returns the normalized value depending on the tpye of the balance change
func (bc Transfer) Value() *big.Int {
	switch bc.Type {
	case Received:
		return new(big.Int).Set(bc.Amount)
	case Spent:
		return new(big.Int).Neg(bc.Amount)
	default:
		return big.NewInt(0)
	}
}

// AccountMovements represents the total change
// made to Account in a certain time range
type AccountMovements struct {
	Address   string
	Transfers []*Transfer
}

// NewAccountMovements creates a new instance of AccountMovement
func NewAccountMovements(address string) *AccountMovements {
	return &AccountMovements{
		Address:   address,
		Transfers: make([]*Transfer, 0),
	}
}

// Sort sorts the changes by block height in ascending order
func (am *AccountMovements) Sort() *AccountMovements {
	sort.Slice(am.Transfers, func(i, j int) bool {
		return am.Transfers[i].BlockHeight < am.Transfers[j].BlockHeight
	})
	return am
}

// Receive adds a transfer as received to the list of changes at the given block height
func (am *AccountMovements) Receive(blockHeight uint64, timestamp uint64, txHash string, amount *big.Int, address string) {
	am.Transfers = append(
		am.Transfers,
		&Transfer{
			Address:     address,
			Amount:      new(big.Int).Set(amount),
			BlockHeight: blockHeight,
			Timestamp:   timestamp,
			TxHash:      txHash,
			Type:        Received,
		},
	)
}

// Spend adds a transfer as spent to the list of changes at the given block height
func (am *AccountMovements) Spend(blockHeight uint64, timestamp uint64, txHash string, amount *big.Int, address string) {
	am.Transfers = append(
		am.Transfers,
		&Transfer{
			Address:     address,
			Amount:      new(big.Int).Set(amount),
			BlockHeight: blockHeight,
			Timestamp:   timestamp,
			TxHash:      txHash,
			Type:        Spent,
		},
	)
}

// AccountAssetsMovedEvent represents a domain event upon AccountMovements
type AccountAssetsMovedEvent struct {
	version    int
	occurredOn time.Time
	subsID     string
	account    string
	ts         []*Transfer
	c          Currency
}

// NewAccountAssetsMovedEvent creates a new instance from AccountMovements
func NewAccountAssetsMovedEvent(subsID string, account string, c Currency, ts []*Transfer) *AccountAssetsMovedEvent {
	return &AccountAssetsMovedEvent{
		version:    1,
		occurredOn: time.Now(),
		subsID:     subsID,
		account:    account,
		c:          c,
		ts:         ts,
	}
}

// Account returns Account property
func (evt *AccountAssetsMovedEvent) Account() string {
	return evt.account
}

// Currency returns the currency property
func (evt *AccountAssetsMovedEvent) Currency() Currency {
	return evt.c
}

// SubscriptionID returns the subsID property
func (evt *AccountAssetsMovedEvent) SubscriptionID() string {
	return evt.subsID
}

// Transfers returns transfers property
func (evt *AccountAssetsMovedEvent) Transfers() []*Transfer {
	return evt.ts
}

// OccurredOn returns event time
func (evt *AccountAssetsMovedEvent) OccurredOn() time.Time {
	return evt.occurredOn
}

// EventVersion returns event version
func (evt *AccountAssetsMovedEvent) EventVersion() int {
	return evt.version
}
