package cryptobot

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
)

// Represents errors related to subscription
var (
	ErrInvalidID = errors.New("invalid identity")
)

// SubscriptionRepository represents common API for subscriptions repository
type SubscriptionRepository interface {
	UnitOfWork
	// NextIdentity returns the next available identity
	NextIdentity(userID string) string
	// Size returns the total number of subscriptions persited in the repository
	Size() int64
	// Get returns the subscription for the given subscription id
	Get(id string) (*Subscription, error)
	// GetAllForUser returns all subscriptions for the given user id
	GetAllForUser(userID string) ([]*Subscription, error)
	// GetAllForCurrency returns all subscriptions for the given currency
	GetAllForCurrency(currencySymbol string) ([]*Subscription, error)
	// Save persists/updates the given subscription
	Save(s *Subscription) error
	// Remove removes the given subscription from the persistance
	Remove(s *Subscription) error
}

// Subscription is a root aggragate
type Subscription struct {
	id                  string
	userID              string
	blockHeight         uint64
	startingBlockHeight uint64
	c                   Currency
	account             string
	totalReceived       *big.Int
	totalSpent          *big.Int
}

// UserIDFrom extracts UserID from SubscriptionID.
// In case of failure, returns empty string.
// subscriptionID = userID + ':' + uuid
func UserIDFrom(subscriptionID string) string {
	s := strings.Split(subscriptionID, ":")

	if s == nil || len(s) < 2 || s[0] == "" {
		return ""
	}

	return s[0]
}

// NewSubscription creates a new subscription
func NewSubscription(
	id string,
	userID string,
	account string,
	c Currency,
	startingBlockHeight uint64,
) (*Subscription, error) {
	if id == "" {
		return nil, ErrInvalidID
	}

	s := &Subscription{
		id:                  id,
		userID:              userID,
		c:                   c,
		account:             account,
		totalReceived:       new(big.Int),
		totalSpent:          new(big.Int),
		blockHeight:         startingBlockHeight,
		startingBlockHeight: startingBlockHeight,
	}

	return s, nil
}

// DeepCopySubscription creates a copy
func DeepCopySubscription(
	id string,
	userID string,
	account string,
	c Currency,
	totalReceived *big.Int,
	totalSpent *big.Int,
	blockHeight uint64,
	staringBlockHeight uint64,
) (*Subscription, error) {
	s, err := NewSubscription(id, userID, account, c, staringBlockHeight)
	if err != nil {
		return nil, err
	}

	s.blockHeight = blockHeight
	s.totalReceived = totalReceived
	s.totalSpent = totalSpent

	return s, nil
}

// Account returns address of this account
func (s *Subscription) Account() string {
	return s.account
}

// TotalReceived returns the total received balance
// since the starting blockheight of the subscription
func (s *Subscription) TotalReceived() *big.Int {
	return s.totalReceived
}

// TotalSpent returns the total spent balance
// since the starting blockheight of the subscription
func (s *Subscription) TotalSpent() *big.Int {
	return s.totalSpent
}

// BlockHeight returns the last block height that balance is updated
func (s *Subscription) BlockHeight() uint64 {
	return s.blockHeight
}

// StartingBlockHeight returns the staring block height
func (s *Subscription) StartingBlockHeight() uint64 {
	return s.startingBlockHeight
}

// Currency returns currency property
func (s *Subscription) Currency() Currency {
	return s.c
}

// ID returns id property
func (s *Subscription) ID() string {
	return s.id
}

// UserID returns userID property
func (s *Subscription) UserID() string {
	return s.userID
}

// ToString returns a string representation for this subscription
func (s *Subscription) ToString() string {
	return fmt.Sprintf(
		"ID: %s\nType: %s\nAsset: %s\nTotalReceived: %s\nTotalSpent: %s\nStarting Block Height: %d\nLast Updated Block Height: %d",
		s.ID(),
		s.Account(),
		s.Currency().Symbol,
		s.TotalReceived().String(),
		s.TotalSpent().String(),
		s.StartingBlockHeight(),
		s.BlockHeight(),
	)
}

// ApplyMovements applies a set of movements to the current state of this account
// Movements in a AccountMovements object should be descending-ordered
// by block height. Otherwise after the first Movement applied, the
// remaining will be ignored.
func (s *Subscription) ApplyMovements(acms *AccountMovements) {
	if acms == nil || acms.Address != s.account {
		log.Printf("account's address(%s) doesn't match with the movement's address(%s), not applying", s.Account(), acms.Address)
		return
	}

	for _, blockHeight := range acms.Blocks {
		if blockHeight <= s.BlockHeight() {
			log.Printf("movement's blockheight(%d) is less than the last updated blockheight(%d), not applying", blockHeight, s.BlockHeight())
			return
		}

		for _, c := range acms.Changes[blockHeight] {
			switch c.Type {
			case ReceivedBalance:
				s.receiveBalance(c.Amount)
				break
			case SpentBalance:
				s.spendBalance(c.Amount)
				break
			default:
			}
		}

		s.setBlockHeight(blockHeight)
	}

	DomainEventPublisherInstance().Publish(
		NewAccountAssetsMovedEvent(s.ID(), s.Currency(), acms))
}

func (s *Subscription) receiveBalance(b *big.Int) {
	s.totalReceived = new(big.Int).Add(s.totalReceived, b)
}

func (s *Subscription) spendBalance(b *big.Int) {
	s.totalSpent = new(big.Int).Add(s.totalSpent, b)
}

func (s *Subscription) setBlockHeight(h uint64) {
	s.blockHeight = h
}
