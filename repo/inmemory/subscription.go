package inmemory

import (
	"errors"

	cryptoBot "github.com/psychoplasma/crypto-balance-bot"
)

var errIndifferentUserID = errors.New("updating UserID field of an existing subscription is not allowed")

// SubscriptionReposititory is an in-memory implementation of SubscriptionReposititory
type SubscriptionReposititory struct {
	subsByUserID  map[string]map[string]*cryptoBot.Subscription
	subsByID      map[string]*cryptoBot.Subscription
	subsActivated map[string]*cryptoBot.Subscription
	size          int
}

// NewSubscriptionReposititory creates a new instance of SubscriptionReposititory
func NewSubscriptionReposititory() *SubscriptionReposititory {
	return &SubscriptionReposititory{
		subsByUserID:  make(map[string]map[string]*cryptoBot.Subscription),
		subsByID:      make(map[string]*cryptoBot.Subscription),
		subsActivated: make(map[string]*cryptoBot.Subscription),
		size:          0,
	}
}

// Size returns the total number of subscriptions persited in the repository
func (r *SubscriptionReposititory) Size() int {
	return r.size
}

// Get returns the subscription for the given subscription id
func (r *SubscriptionReposititory) Get(id string) (*cryptoBot.Subscription, error) {
	return r.subsByID[id], nil
}

// GetAllForUser returns all subscriptions for the given user id
func (r *SubscriptionReposititory) GetAllForUser(userID string) ([]*cryptoBot.Subscription, error) {
	subs := make([]*cryptoBot.Subscription, 0)
	for _, s := range r.subsByUserID[userID] {
		subs = append(subs, s)
	}
	return subs, nil
}

// GetAllActivated returns all activated subscriptions
func (r *SubscriptionReposititory) GetAllActivated() ([]*cryptoBot.Subscription, error) {
	subs := make([]*cryptoBot.Subscription, 0)
	for _, s := range r.subsActivated {
		subs = append(subs, s)
	}
	return subs, nil
}

// Add persists/updates the given subscription
func (r *SubscriptionReposititory) Add(s *cryptoBot.Subscription) error {
	// Do not allow to update UserID of an existing subscription
	if r.subsByID[s.ID] != nil && r.subsByID[s.ID].UserID != s.UserID {
		return errIndifferentUserID
	}

	// Increment the size if the item doesn't exit upon persistance
	if r.subsByID[s.ID] == nil {
		// We're caching the size because everytime calling
		// len() would have an overhead for a non-local slice
		r.size++
	}
	r.subsByID[s.ID] = s

	if r.subsByUserID[s.UserID] == nil {
		r.subsByUserID[s.UserID] = make(map[string]*cryptoBot.Subscription)
	}
	r.subsByUserID[s.UserID][s.ID] = s

	if s.Activated {
		r.subsActivated[s.ID] = s
	} else {
		delete(r.subsActivated, s.ID)
	}

	return nil
}

// Remove removes the given subscription from the persistance
func (r *SubscriptionReposititory) Remove(s *cryptoBot.Subscription) error {
	// Decrement the size if the item exists upon removal
	if r.subsByID[s.ID] != nil {
		r.size--
	}

	delete(r.subsByID, s.ID)
	delete(r.subsByUserID[s.UserID], s.ID)
	delete(r.subsActivated, s.ID)

	return nil
}
