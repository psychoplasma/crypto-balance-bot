package application

import (
	domain "github.com/psychoplasma/crypto-balance-bot"
)

// SubscriptionApplication exposes application services for subscription entity
type SubscriptionApplication struct {
	r domain.SubscriptionRepository
}

// NewSubscriptionApplication fatory function
func NewSubscriptionApplication(r domain.SubscriptionRepository) *SubscriptionApplication {
	return &SubscriptionApplication{
		r: r,
	}
}

// SubscribeForValue creates a new value-based subscription and activates it
func (sa *SubscriptionApplication) SubscribeForValue(userID string, name string, c domain.Currency, addrDesc string, against domain.Currency) error {
	s, err := domain.NewSubscription(sa.r.NextIdentity(), userID, name, domain.ValueSubscription, against)
	if err != nil {
		return err
	}
	s.AddAccount(c, addrDesc)
	s.Activate()

	return sa.r.Add(s)
}

// SubscribeForMovement creates a new movement-based subscription and activates it
func (sa *SubscriptionApplication) SubscribeForMovement(userID string, name string, c domain.Currency, addrDesc string) error {
	s, err := domain.NewSubscription(sa.r.NextIdentity(), userID, name, domain.MovementSubscription, domain.Currency{})
	if err != nil {
		return err
	}
	s.AddAccount(c, addrDesc)
	s.Activate()

	return sa.r.Add(s)
}

// Unsubscribe removes the given subscription
func (sa *SubscriptionApplication) Unsubscribe(subscriptionID string) error {
	s, err := sa.r.Get(subscriptionID)
	if err != nil {
		return err
	}
	return sa.r.Remove(s)
}

// UnsubscribeAll removes all subscription belogs to the given user
func (sa *SubscriptionApplication) UnsubscribeAll(userID string) error {
	subs, err := sa.r.GetAllForUser(userID)
	if err != nil {
		return err
	}

	for _, s := range subs {
		if err := sa.r.Remove(s); err != nil {
			return err
		}
	}

	return nil
}

// ActivateSubscription activates the given subscription
func (sa *SubscriptionApplication) ActivateSubscription(subscriptionID string) error {
	s, err := sa.r.Get(subscriptionID)
	if err != nil {
		return err
	}
	s.Activate()
	sa.r.Add(s)

	return nil
}

// DeactivateSubscription deactivates the given subscription
func (sa *SubscriptionApplication) DeactivateSubscription(subscriptionID string) error {
	s, err := sa.r.Get(subscriptionID)
	if err != nil {
		return err
	}
	s.Deactivate()
	sa.r.Add(s)

	return nil
}

// GetSubscription returns the details of the given subscription
func (sa *SubscriptionApplication) GetSubscription(id string) (*domain.Subscription, error) {
	return sa.r.Get(id)
}

// GetSubscriptionsForUser returns the details of all subscriptions for the given user
func (sa *SubscriptionApplication) GetSubscriptionsForUser(userID string) ([]*domain.Subscription, error) {
	return sa.r.GetAllForUser(userID)
}

// GetActiveSubscriptions returns the all active subscriptions
func (sa *SubscriptionApplication) GetActiveSubscriptions() ([]*domain.Subscription, error) {
	return sa.r.GetAllActivated()
}
