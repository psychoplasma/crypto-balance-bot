package application

import (
	"errors"
	"fmt"

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

// Subscribe creates a new subscription
func (sa *SubscriptionApplication) Subscribe(userID string, currencySymbol string, account string) error {
	if err := sa.r.Begin(); err != nil {
		return err
	}

	s, err := sa.subscribe(userID, currencySymbol, account)
	if err != nil {
		return sa.returnError(err)
	}

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

// AddAmountFilter adds a new amount filter
func (sa *SubscriptionApplication) AddAmountFilter(subsID string, amount string, must bool) error {
	if err := sa.r.Begin(); err != nil {
		return err
	}

	s, err := sa.r.Get(subsID)
	if err != nil {
		return sa.returnError(err)
	}

	f, err := domain.NewAmountFilter(amount, must)
	if err != nil {
		return sa.returnError(err)
	}
	s.AddFilter(f)

	if err := sa.r.Save(s); err != nil {
		return sa.returnError(err)
	}

	sa.r.Success()

	return nil
}

// AddAddressOnFilter adds a new address-off filter
func (sa *SubscriptionApplication) AddAddressOnFilter(subsID string, address string, must bool) error {
	if err := sa.r.Begin(); err != nil {
		return err
	}

	s, err := sa.r.Get(subsID)
	if err != nil {
		return sa.returnError(err)
	}

	f, err := domain.NewAddressOnFilter(address, must)
	if err != nil {
		return sa.returnError(err)
	}
	s.AddFilter(f)

	if err := sa.r.Save(s); err != nil {
		return sa.returnError(err)
	}

	sa.r.Success()

	return nil
}

// AddAddressOffFilter adds a new address-on filter
func (sa *SubscriptionApplication) AddAddressOffFilter(subsID string, address string, must bool) error {
	if err := sa.r.Begin(); err != nil {
		return err
	}

	s, err := sa.r.Get(subsID)
	if err != nil {
		return sa.returnError(err)
	}

	f, err := domain.NewAddressOffFilter(address, must)
	if err != nil {
		return sa.returnError(err)
	}
	s.AddFilter(f)

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

// GetSubscriptionsForCurrency returns all subscriptions for a given currency that are updated before the given blockheight
func (sa *SubscriptionApplication) GetSubscriptionsForCurrency(currencySymbol string, updatedBefore uint64) ([]*domain.Subscription, error) {
	if err := sa.r.Begin(); err != nil {
		return nil, err
	}

	cs, exist := services.CurrencyServiceFactory[currencySymbol]
	if !exist {
		return nil, errInexistentCurrency
	}

	cs.GetLatestBlockHeight()

	subs, err := sa.r.GetAllForCurrency(currencySymbol, updatedBefore)
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

func (sa *SubscriptionApplication) subscribe(userID string, currencySymbol string, account string) (*domain.Subscription, error) {
	c, exist := services.CurrencyFactory[currencySymbol]
	if !exist {
		return nil, errInexistentCurrency
	}

	cs, exist := services.CurrencyServiceFactory[currencySymbol]
	if !exist {
		return nil, errInexistentCurrency
	}

	bh, err := cs.GetLatestBlockHeight()
	if err != nil {
		return nil, err
	}

	s, err := domain.NewSubscription(
		sa.r.NextIdentity(userID),
		userID,
		account,
		c,
		bh,
	)
	if err != nil {
		return nil, err
	}

	return s, nil
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
