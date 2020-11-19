package cryptobot

import (
	"errors"
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

// SubscriptionRepository repostitory for subscriptions
type SubscriptionRepository interface {
	NextIdentity() string
	Size() int
	Get(id string) (*Subscription, error)
	GetAllForUser(userID string) ([]*Subscription, error)
	GetAllActivated() ([]*Subscription, error)
	Add(s *Subscription) error
	Remove(s *Subscription) error
}

// Subscription is a root aggragate
type Subscription struct {
	ID              string
	UserID          string
	Name            string
	Type            SubscriptionType
	Activated       bool
	AgainstCurrency Currency
	Accounts        map[string]*Account
}

// NewSubscription creates a new subscription
func NewSubscription(id string, userID string, name string, stype SubscriptionType, against Currency) (*Subscription, error) {
	if id == "" {
		return nil, ErrInvalidID
	}

	if isSubsctiptionTypeValid(stype) {
		return nil, ErrInvalidSubscriptionType
	}

	s := &Subscription{
		ID:              id,
		UserID:          userID,
		Name:            name,
		Type:            stype,
		AgainstCurrency: against,
		Accounts:        make(map[string]*Account),
	}

	return s, nil
}

func (s *Subscription) AddAccount(c Currency, address string) {
	if s.Accounts[address] != nil {
		return
	}

	s.Accounts[address] = NewAccount(c, address)
}

// Activate activates the subscription. User will start getting notifications about this subscription
func (s *Subscription) Activate() {
	if s.Activated {
		return
	}
	s.Activated = true
}

// Deactivate deactivates the subscription. User will stop getting notifications about this subscription
func (s *Subscription) Deactivate() {
	if !s.Activated {
		return
	}
	s.Activated = false
}

func isSubsctiptionTypeValid(stype SubscriptionType) bool {
	return stype != ValueSubscription &&
		stype != MovementSubscription
}
