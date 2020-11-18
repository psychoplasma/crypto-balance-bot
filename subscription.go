package cryptobot

import (
	"errors"

	"github.com/google/uuid"
)

// Represents errors related to subscription
var (
	ErrInvalidSubscriptionType = errors.New("invalid subscription type")
)

// SubscriptionType enumarations
type SubscriptionType string

// Values for SubscriptionType
const (
	Value    = SubscriptionType("value")
	Movement = SubscriptionType("movement")
)

// SubscriptionRepository repostitory for subscriptions
type SubscriptionRepository interface {
	Size() int
	Get(id string) (*Subscription, error)
	GetAllForUser(userID string) ([]*Subscription, error)
	GetAllActivated() ([]*Subscription, error)
	Add(s *Subscription) error
	Remove(s *Subscription) error
}

// Subscription implements un/subscription logic for crypobot
type Subscription struct {
	ID              string
	UserID          string
	Name            string
	Type            SubscriptionType
	Activated       bool
	AgainstCurrency Currency
	Account         *Account
}

// ValueSubscription creates a new value-based subscription
func ValueSubscription(userID string, name string, c Currency, addrDesc string, against Currency) (*Subscription, error) {
	return NewSubscription(userID, name, SubscriptionType(Movement), c, addrDesc, against)
}

// MovementSubscription creates a new movement-based subscription
func MovementSubscription(userID string, name string, c Currency, addrDesc string) (*Subscription, error) {
	return NewSubscription(userID, name, SubscriptionType(Value), c, addrDesc, Currency{})
}

// NewSubscription creates a new subscription
func NewSubscription(userID string, name string, stype SubscriptionType, c Currency, addrDesc string, against Currency) (*Subscription, error) {
	if stype != Value && stype != Movement {
		return nil, ErrInvalidSubscriptionType
	}

	id := uuid.New()
	s := &Subscription{
		Name:   name,
		UserID: userID,
		ID:     id.String(),
		Type:   stype,
	}

	a, err := NewAccount(c, addrDesc)
	if err != nil {
		return nil, err
	}

	s.Account = a
	s.AgainstCurrency = against

	return s, nil
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
