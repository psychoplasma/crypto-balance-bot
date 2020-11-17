package application

import (
	cryptoBot "github.com/psychoplasma/crypto-balance-bot"
)

// SubscriptionService exposes application services for subscription entity
type SubscriptionService struct {
	r cryptoBot.SubscriptionRepository
}

// SubscriptionService fatory function
func SubscriptionService(r cryptoBot.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		r: r,
	}
}

// SubscribeForValue creates a new value-based subscription and activates it
func (ss *SubscriptionService) SubscribeForValue(userID string, name string, c cryptoBot.Currency, addrDesc string, against cryptoBot.Currency) error {
	s, err := cryptoBot.ValueSubscription(userID, name, c, addrDesc, against)
	if err != nil {
		return err
	}
	s.Activate()

	return ss.r.Add(s)
}

// SubscribeForMovement creates a new movement-based subscription and activates it
func (ss *SubscriptionService) SubscribeForMovement(userID string, name string, c cryptoBot.Currency, addrDesc string, against cryptoBot.Currency) error {
	s, err := cryptoBot.MovementSubscription(userID, name, c, addrDesc, against)
	if err != nil {
		return err
	}
	s.Activate()

	return ss.r.Add(s)
}

// Unsubscribe removes the given subscription
func (ss *SubscriptionService) Unsubscribe(subscriptionID string) error {
	return ss.r.Remove(subscriptionID)
}

// UnsubscribeAll removes all subscription belogs to the given user
func (ss *SubscriptionService) UnsubscribeAll(userID string) error {
	subs, err := ss.r.GetAllForUser(userID)
	if err != nil {
		return err
	}

	for _, s := range subs {
		if err := ss.r.Remove(s.ID); err != nil {
			return err
		}
	}

	return nil
}

// ActivateSubscription activates the given subscription
func (ss *SubscriptionService) ActivateSubscription(subscriptionID string) error {
	s, err := ss.r.Get(subscriptionID)
	if err != nil {
		return err
	}
	s.Activate()

	return nil
}

// DeactivateSubscription deactivates the given subscription
func (ss *SubscriptionService) DeactivateSubscription(subscriptionID string) error {
	s, err := ss.r.Get(subscriptionID)
	if err != nil {
		return err
	}
	s.Deactivate()

	return nil
}

// GetSubscription returns the details of the given subscription
func (ss *SubscriptionService) GetSubscription(userID string) (*cryptoBot.Subscription, error) {
	return ss.r.Get(userID)
}

// GetSubscriptionsForUser returns the details of all subscriptions for the given user
func (ss *SubscriptionService) GetSubscriptionsForUser(userID string) ([]*cryptoBot.Subscription, error) {
	subs, err := ss.r.GetAllForUser(userID)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

// GetActiveSubscriptions returns the active subscriptions
func (ss *SubscriptionService) GetActiveSubscriptions(userID string) ([]*cryptoBot.Subscription, error) {
	// FIXME:
	subs, err := ss.r.GetAllForUser(userID)
	if err != nil {
		return nil, err
	}

	return subs, nil
}
