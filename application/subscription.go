package application

import (
	"errors"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/services"
)

var (
	errInexistentCurrency = errors.New("inexistent currency")
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
func (sa *SubscriptionApplication) SubscribeForValue(userID string, name string, currencySymbol string, againstCurrencySymbol string, addrDescs []string) error {
	c, e := services.CurrencyFactory[currencySymbol]
	if !e {
		return errInexistentCurrency
	}

	ac, e := services.CurrencyFactory[currencySymbol]
	if !e {
		return errInexistentCurrency
	}

	s, err := domain.NewSubscription(sa.r.NextIdentity(), userID, name, domain.ValueSubscription, *ac)
	if err != nil {
		return err
	}

	for _, addr := range addrDescs {
		s.AddAccount(*c, addr)
	}
	s.Activate()

	return sa.r.Save(s)
}

// SubscribeForMovement creates a new movement-based subscription and activates it
func (sa *SubscriptionApplication) SubscribeForMovement(userID string, name string, currencySymbol string, addrDescs []string) error {
	c, e := services.CurrencyFactory[currencySymbol]
	if !e {
		return errInexistentCurrency
	}

	s, err := domain.NewSubscription(sa.r.NextIdentity(), userID, name, domain.MovementSubscription, domain.Currency{})
	if err != nil {
		return err
	}

	for _, addr := range addrDescs {
		s.AddAccount(*c, addr)
	}
	s.Activate()

	return sa.r.Save(s)
}

// Unsubscribe removes the given subscription
func (sa *SubscriptionApplication) Unsubscribe(subscriptionID string) error {
	s, err := sa.r.Get(subscriptionID)
	if err != nil {
		return err
	}
	return sa.r.Remove(s)
}

// UnsubscribeAllForUser removes all subscription belogs to the given user
func (sa *SubscriptionApplication) UnsubscribeAllForUser(userID string) error {
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
	sa.r.Save(s)

	return nil
}

// DeactivateSubscription deactivates the given subscription
func (sa *SubscriptionApplication) DeactivateSubscription(subscriptionID string) error {
	s, err := sa.r.Get(subscriptionID)
	if err != nil {
		return err
	}
	s.Deactivate()
	sa.r.Save(s)

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
