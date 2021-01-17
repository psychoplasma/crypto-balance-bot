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
	ErrInvalidSubscriptionType = errors.New("invalid subscription type")
	ErrInvalidID               = errors.New("invalid identity")
)

// SubscriptionType enumarations
type SubscriptionType string

// Values for SubscriptionType
const (
	ValueSubscription    = SubscriptionType("value")
	MovementSubscription = SubscriptionType("movement")
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
	// GetAllActivatedMovements returns all activated movement subscriptions
	GetAllActivatedMovements() ([]*Subscription, error)
	// GetAllActivatedValues returns all activated value subscriptions
	GetAllActivatedValues() ([]*Subscription, error)
	// Save persists/updates the given subscription
	Save(s *Subscription) error
	// Remove removes the given subscription from the persistance
	Remove(s *Subscription) error
}

// Subscription is a root aggragate
type Subscription struct {
	id          string
	userID      string
	stype       SubscriptionType
	activated   bool
	c           Currency
	ac          Currency
	account     string
	balance     *big.Int
	blockHeight int
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
	stype SubscriptionType,
	account string,
	c Currency,
	against Currency,
) (*Subscription, error) {
	if id == "" {
		return nil, ErrInvalidID
	}

	if isSubsctiptionTypeValid(stype) {
		return nil, ErrInvalidSubscriptionType
	}

	s := &Subscription{
		id:      id,
		userID:  userID,
		stype:   stype,
		c:       c,
		ac:      against,
		account: account,
		balance: new(big.Int),
	}

	return s, nil
}

// DeepCopySubscription creates a copy
func DeepCopySubscription(
	id string,
	userID string,
	stype SubscriptionType,
	activated bool,
	account string,
	c Currency,
	against Currency,
	balance *big.Int,
	blockHeight int,
) (*Subscription, error) {
	s, err := NewSubscription(id, userID, stype, account, c, against)
	if err != nil {
		return nil, err
	}

	s.setBalance(balance)
	s.setBlockHeight(blockHeight)

	if activated {
		s.Activate()
	}

	return s, nil
}

// Account returns address of this account
func (s *Subscription) Account() string {
	return s.account
}

// Balance returns the last checked balance
func (s *Subscription) Balance() *big.Int {
	return s.balance
}

// BlockHeight returns the last block height that balance is updated
func (s *Subscription) BlockHeight() int {
	return s.blockHeight
}

// AgainstCurrency returns against currency property
func (s *Subscription) AgainstCurrency() Currency {
	return s.ac
}

// Currency returns currency property
func (s *Subscription) Currency() Currency {
	return s.c
}

// ID returns id property
func (s *Subscription) ID() string {
	return s.id
}

// IsActivated returns activated property
func (s *Subscription) IsActivated() bool {
	return s.activated
}

// UserID returns userID property
func (s *Subscription) UserID() string {
	return s.userID
}

// Type returns stype property
func (s *Subscription) Type() SubscriptionType {
	return s.stype
}

// ToString returns a string representation for this subscription
func (s *Subscription) ToString() string {
	status := ""
	if s.IsActivated() {
		status = "Active"
	} else {
		status = "Deactive"
	}

	log.Printf("%+v", s)

	return fmt.Sprintf(
		"ID: %s\nType: %s\nAsset: %s\nStatus: %s\nBalance: %s\nLast Updated Block Height: %d",
		s.ID(),
		s.Account(),
		s.Currency().Symbol,
		status,
		s.Balance().String(),
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
			s.setBalance(new(big.Int).Add(s.Balance(), c.Amount))
		}

		s.setBlockHeight(blockHeight)
	}

	DomainEventPublisherInstance().Publish(
		NewAccountAssetsMovedEvent(s.ID(), s.Currency(), acms))
}

// Activate activates the subscription. User will start getting notifications about this subscription
func (s *Subscription) Activate() {
	if s.activated {
		return
	}
	s.activated = true
}

// Deactivate deactivates the subscription. User will stop getting notifications about this subscription
func (s *Subscription) Deactivate() {
	if !s.activated {
		return
	}
	s.activated = false
}

func (s *Subscription) setBalance(b *big.Int) {
	s.balance = b
}

func (s *Subscription) setBlockHeight(h int) {
	s.blockHeight = h
}

func isSubsctiptionTypeValid(stype SubscriptionType) bool {
	return stype != ValueSubscription &&
		stype != MovementSubscription
}
