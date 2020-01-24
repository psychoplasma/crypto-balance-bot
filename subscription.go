package cryptobot

import (
	"errors"
)

const (
	value = iota
	movement
)

var errSubscriptionType = errors.New("wrong subscription type")

type SubscriptionType int

type Subscription struct {
	Type            SubscriptionType
	AgainstCurrency string
	Account         *Account
}

// Subscribe subscribes for value change or account movement events for the given account
func (s *Subscription) Subscribe(stype SubscriptionType) error {
	switch stype {
	case value:
		return s.subscribeForValue()
	case movement:
		return s.subscribeForMovement()
	default:
		return errSubscriptionType
	}
}

// Unsubscribe unsubscribes for value change or account movement events for the given account
func (s *Subscription) Unsubscribe(stype SubscriptionType) error {
	switch stype {
	case value:
		return s.unsubscribeForValue(a)
	case movement:
		return s.unsubscribeForMovement(a)
	default:
		return errSubscriptionType
	}
}

func (s *Subscription) subscribeForValue() error {
	return nil
}

func (s *Subscription) subscribeForMovement() error {
	return nil
}

func (s *Subscription) unsubscribeForValue() error {
	return nil
}

func (s *Subscription) unsubscribeForMovement() error {
	return nil
}
