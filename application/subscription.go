package application

import (
	"errors"
	"fmt"
	"log"

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

// NewSubscriptionApplication factory function
func NewSubscriptionApplication(repo domain.SubscriptionRepository) *SubscriptionApplication {
	return &SubscriptionApplication{
		r: repo,
	}
}

// SubscribeForValue creates a new value-based subscription and activates it
func (sa *SubscriptionApplication) SubscribeForValue(userID string, currencySymbol string, againstCurrencySymbol string, account string) error {
	if err := sa.r.Begin(); err != nil {
		return err
	}

	c, exist := services.CurrencyFactory[currencySymbol]
	if !exist {
		return sa.returnError(errInexistentCurrency)
	}

	ac, exist := services.CurrencyFactory[againstCurrencySymbol]
	if !exist {
		return sa.returnError(errInexistentCurrency)
	}

	cs, exist := services.CurrencyServiceFactory[currencySymbol]
	if !exist {
		return sa.returnError(errInexistentCurrency)
	}

	bh, err := cs.GetLatestBlockHeight()
	if err != nil {
		return sa.returnError(err)
	}

	s, err := domain.NewSubscription(
		sa.r.NextIdentity(userID),
		userID,
		domain.ValueSubscription,
		account,
		c,
		ac,
		bh,
	)
	if err != nil {
		return sa.returnError(err)
	}

	s.Activate()

	if err := sa.r.Save(s); err != nil {
		return sa.returnError(err)
	}

	sa.r.Success()

	return nil
}

// SubscribeForMovement creates a new movement-based subscription and activates it
func (sa *SubscriptionApplication) SubscribeForMovement(userID string, currencySymbol string, account string) error {
	if err := sa.r.Begin(); err != nil {
		return err
	}

	c, exist := services.CurrencyFactory[currencySymbol]
	if !exist {
		sa.r.Fail()
		return errInexistentCurrency
	}

	cs, exist := services.CurrencyServiceFactory[currencySymbol]
	if !exist {
		return sa.returnError(errInexistentCurrency)
	}

	bh, err := cs.GetLatestBlockHeight()
	if err != nil {
		return sa.returnError(err)
	}

	s, err := domain.NewSubscription(
		sa.r.NextIdentity(userID),
		userID,
		domain.MovementSubscription,
		account,
		c,
		domain.Currency{},
		bh,
	)
	if err != nil {
		return sa.returnError(err)
	}

	s.Activate()

	if err := sa.r.Save(s); err != nil {
		return sa.returnError(err)
	}

	sa.r.Success()

	return nil
}

// Unsubscribe removes the given subscription
func (sa *SubscriptionApplication) Unsubscribe(subscriptionID string) error {
	if err := sa.r.Begin(); err != nil {
		return err
	}

	s, err := sa.r.Get(subscriptionID)
	if err != nil {
		return sa.returnError(err)
	}

	if err := sa.r.Remove(s); err != nil {
		return sa.returnError(err)
	}

	sa.r.Success()

	return nil
}

// UnsubscribeAllForUser removes all subscription belogs to the given user
func (sa *SubscriptionApplication) UnsubscribeAllForUser(userID string) error {
	if err := sa.r.Begin(); err != nil {
		return err
	}

	subs, err := sa.r.GetAllForUser(userID)
	if err != nil {
		return sa.returnError(err)
	}

	for _, s := range subs {
		if err := sa.r.Remove(s); err != nil {
			return sa.returnError(err)
		}
	}

	sa.r.Success()

	return nil
}

// ActivateSubscription activates the given subscription
func (sa *SubscriptionApplication) ActivateSubscription(subscriptionID string) error {
	if err := sa.r.Begin(); err != nil {
		return err
	}

	s, err := sa.r.Get(subscriptionID)
	if err != nil {
		return sa.returnError(err)
	}

	s.Activate()

	if err := sa.r.Save(s); err != nil {
		return sa.returnError(err)
	}

	sa.r.Success()

	return nil
}

// DeactivateSubscription deactivates the given subscription
func (sa *SubscriptionApplication) DeactivateSubscription(subscriptionID string) error {
	if err := sa.r.Begin(); err != nil {
		return err
	}

	s, err := sa.r.Get(subscriptionID)
	if err != nil {
		return sa.returnError(err)
	}

	s.Deactivate()

	if err := sa.r.Save(s); err != nil {
		return sa.returnError(err)
	}

	sa.r.Success()

	return nil
}

// GetSubscription returns the details of the given subscription
func (sa *SubscriptionApplication) GetSubscription(id string) (*domain.Subscription, error) {
	if err := sa.r.Begin(); err != nil {
		return nil, err
	}

	s, err := sa.r.Get(id)
	if err != nil {
		return nil, sa.returnError(err)
	}

	sa.r.Success()

	return s, nil
}

// GetSubscriptionsForUser returns the details of all subscriptions for the given user
func (sa *SubscriptionApplication) GetSubscriptionsForUser(userID string) ([]*domain.Subscription, error) {
	if err := sa.r.Begin(); err != nil {
		return nil, err
	}

	subs, err := sa.r.GetAllForUser(userID)
	if err != nil {
		return nil, sa.returnError(err)
	}

	sa.r.Success()

	return subs, nil
}

// GetAllActivatedMovements returns all activated movement subscriptions
func (sa *SubscriptionApplication) GetAllActivatedMovements() ([]*domain.Subscription, error) {
	if err := sa.r.Begin(); err != nil {
		return nil, err
	}

	subs, err := sa.r.GetAllActivatedMovements()
	if err != nil {
		return nil, sa.returnError(err)
	}

	sa.r.Success()

	return subs, nil
}

// CheckAndApplyAccountMovements checks whether there is any movement
// for the given account and if there is, applies them to the account.
func (sa *SubscriptionApplication) CheckAndApplyAccountMovements(s *domain.Subscription) error {
	if err := sa.r.Begin(); err != nil {
		return err
	}

	if err := sa.checkAndApplyAccountMovements(s); err != nil {
		return sa.returnError(err)
	}

	sa.r.Success()

	return nil
}

// CheckAndApplyAccountMovementsForAllActiveSubscriptions checks and applies
// account movements for all active movement subscriptions
func (sa *SubscriptionApplication) CheckAndApplyAccountMovementsForAllActiveSubscriptions() error {
	if err := sa.r.Begin(); err != nil {
		return err
	}

	subs, err := sa.r.GetAllActivatedMovements()
	if err != nil {
		return sa.returnError(err)
	}

	for _, s := range subs {
		if err := sa.checkAndApplyAccountMovements(s); err != nil {
			log.Printf("error while checking account movements for %s : %s", s.Account(), err.Error())
		}
	}

	sa.r.Success()

	return nil
}

func (sa *SubscriptionApplication) checkAndApplyAccountMovements(s *domain.Subscription) error {
	if s == nil {
		return fmt.Errorf("nil subscription")
	}

	cs, exist := services.CurrencyServiceFactory[s.Currency().Symbol]
	if !exist {
		return fmt.Errorf("no currency service found for %s", s.Currency().Symbol)
	}

	acm, err := cs.GetAccountMovements(s.Account(), s.BlockHeight()+1)
	if err != nil {
		return err
	}

	s.ApplyMovements(acm.Sort())

	return sa.r.Save(s)
}

func (sa *SubscriptionApplication) returnError(err error) error {
	sa.r.Fail()
	return err
}
