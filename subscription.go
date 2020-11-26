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
	GetAllActivatedMovements() ([]*Subscription, error)
	GetAllActivatedValues() ([]*Subscription, error)
	Save(s *Subscription) error
	Remove(s *Subscription) error
}

// Subscription is a root aggragate
type Subscription struct {
	id        string
	userID    string
	name      string
	stype     SubscriptionType
	activated bool
	ac        Currency
	accs      []*Account
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
		id:     id,
		userID: userID,
		name:   name,
		stype:  stype,
		ac:     against,
		accs:   make([]*Account, 0),
	}

	return s, nil
}

// ID returns id property
func (s *Subscription) ID() string {
	return s.id
}

// UserID returns userID property
func (s *Subscription) UserID() string {
	return s.userID
}

// Name returns name property
func (s *Subscription) Name() string {
	return s.name
}

// Type returns stype property
func (s *Subscription) Type() SubscriptionType {
	return s.stype
}

// IsActivated returns activated property
func (s *Subscription) IsActivated() bool {
	return s.activated
}

// Accounts returns accounts property
func (s *Subscription) Accounts() []*Account {
	return s.accs
}

// AddAccount adds a new account to this subscriptions. Duplicates will be overwritten
func (s *Subscription) AddAccount(c Currency, address string) {
	for _, a := range s.Accounts() {
		if a.Address() == address {
			return
		}
	}

	s.accs = append(s.accs, NewAccount(c, address))
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

func isSubsctiptionTypeValid(stype SubscriptionType) bool {
	return stype != ValueSubscription &&
		stype != MovementSubscription
}
