package cryptobot

import (
	"log"
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

// OccurredOn returns event time
func (evt *AccountAssetsMovedEvent) OccurredOn() time.Time {
	return evt.occurredOn
}

// EventVersion returns event version
func (evt *AccountAssetsMovedEvent) EventVersion() int {
	return evt.version
}

// Account in a value object (even it seems like an entity).
// Only balance and index properties change over time
// but they are actually for tracking purposes which is not
// really Address's property.
type Account struct {
	subsID      string
	address     string
	balance     *big.Int
	currency    Currency
	blockHeight int
}

// NewAccount creates a new instance of Account with the given address
func NewAccount(subsID string, address string, c Currency) *Account {
	return &Account{
		address:     address,
		balance:     big.NewInt(0),
		currency:    c,
		blockHeight: -1,
		subsID:      subsID,
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

// SubscriptionID returns the subscription ID of which this account subscribed
func (a *Account) SubscriptionID() string {
	return a.subsID
}

// Apply applies a movement to the current state of this account
// Movements in a AccountMovements object should be descending-ordered
// by block height. Otherwise after the first Movement applied, the
// remaining will be ignored.
func (a *Account) Apply(acms *AccountMovements) {
	if acms == nil || acms.Address != a.address {
		log.Printf("account's address(%s) doesn't match with the movement's address(%s), not applying", a.address, acms.Address)
		return
	}

	for _, blockHeight := range acms.Blocks {
		if blockHeight <= a.blockHeight {
			log.Printf("movement's blockheight(%d) is less than the last updated blockheight(%d), not applying", blockHeight, a.blockHeight)
			return
		}

		for _, c := range acms.Changes[blockHeight] {
			a.balance = a.balance.Add(a.balance, c.Amount)
		}

		a.blockHeight = blockHeight
	}

	DomainEventPublisherInstance().
		Publish(NewAccountAssetsMovedEvent(a.SubscriptionID(), a.Currency(), acms))
}
