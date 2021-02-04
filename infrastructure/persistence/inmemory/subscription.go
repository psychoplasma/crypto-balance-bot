package inmemory

import (
	"errors"

	"github.com/google/uuid"
	domain "github.com/psychoplasma/crypto-balance-bot"
)

var errIndifferentUserID = errors.New("updating UserID field of an existing subscription is not allowed")

// SubscriptionRepository is an in-memory implementation of SubscriptionRepository
type SubscriptionRepository struct {
	subsByUserID  map[string]map[string]*domain.Subscription
	subsByID      map[string]*domain.Subscription
	subsMovements map[string]*domain.Subscription
	subsValues    map[string]*domain.Subscription
	size          int
}

// NewSubscriptionRepository creates a new instance of SubscriptionRepository
func NewSubscriptionRepository() *SubscriptionRepository {
	return &SubscriptionRepository{
		subsByUserID:  make(map[string]map[string]*domain.Subscription),
		subsByID:      make(map[string]*domain.Subscription),
		subsMovements: make(map[string]*domain.Subscription),
		subsValues:    make(map[string]*domain.Subscription),
		size:          0,
	}
}

// Begin starts a new unit for a work to be done on repository
func (r *SubscriptionRepository) Begin() error {
	return nil
}

// Fail rollbacks repository to the state before this work
func (r *SubscriptionRepository) Fail() {}

// Success finalizes the work done on repository
func (r *SubscriptionRepository) Success() {}

// NextIdentity returns the next available identity
func (r *SubscriptionRepository) NextIdentity(userID string) string {
	return userID + ":" + uuid.New().String()
}

// Size returns the total number of subscriptions persited in the repository
func (r *SubscriptionRepository) Size() int64 {
	return int64(r.size)
}

// Get returns the subscription for the given subscription id
func (r *SubscriptionRepository) Get(id string) (*domain.Subscription, error) {
	return r.subsByID[id], nil
}

// GetAllForUser returns all subscriptions for the given user id
func (r *SubscriptionRepository) GetAllForUser(userID string) ([]*domain.Subscription, error) {
	subs := make([]*domain.Subscription, 0)
	for k := range r.subsByUserID[userID] {
		subs = append(subs, r.subsByUserID[userID][k])
	}
	return subs, nil
}

// GetAllMovements returns all movement subscriptions
func (r *SubscriptionRepository) GetAllMovements() ([]*domain.Subscription, error) {
	subs := make([]*domain.Subscription, 0)
	for k := range r.subsMovements {
		subs = append(subs, r.subsMovements[k])
	}
	return subs, nil
}

// GetAllValues returns all value subscriptions
func (r *SubscriptionRepository) GetAllValues() ([]*domain.Subscription, error) {
	subs := make([]*domain.Subscription, 0)
	for k := range r.subsValues {
		subs = append(subs, r.subsValues[k])
	}
	return subs, nil
}

// Save persists/updates the given subscription
func (r *SubscriptionRepository) Save(s *domain.Subscription) error {
	// Do not allow to update UserID of an existing subscription
	if r.subsByID[s.ID()] != nil && r.subsByID[s.ID()].UserID() != s.UserID() {
		return errIndifferentUserID
	}

	// Increment the size if the item doesn't exit upon persistance
	if r.subsByID[s.ID()] == nil {
		// We're caching the size because everytime calling
		// len() would have an overhead for a non-local slice
		r.size++
	}
	r.subsByID[s.ID()] = s

	if r.subsByUserID[s.UserID()] == nil {
		r.subsByUserID[s.UserID()] = make(map[string]*domain.Subscription)
	}
	r.subsByUserID[s.UserID()][s.ID()] = s

	if s.Type() == domain.MovementSubscription {
		r.subsMovements[s.ID()] = s
	}

	if s.Type() == domain.ValueSubscription {
		r.subsValues[s.ID()] = s
	}

	return nil
}

// Remove removes the given subscription from the persistance
func (r *SubscriptionRepository) Remove(s *domain.Subscription) error {
	// Decrement the size if the item exists upon removal
	if r.subsByID[s.ID()] != nil {
		r.size--
	}

	delete(r.subsByID, s.ID())
	delete(r.subsByUserID[s.UserID()], s.ID())
	delete(r.subsMovements, s.ID())
	delete(r.subsValues, s.ID())

	return nil
}
