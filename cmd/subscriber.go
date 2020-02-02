package main

import (
	"errors"
	"fmt"

	cryptobot "github.com/psychoplasma/crypto-balance-bot"
)

// Represents errors related to subscription
var (
	ErrSubscriptionType = errors.New("wrong subscription type")
	ErrNotExists        = errors.New("subscription does not exist")
	ErrAlreadyExists    = errors.New("already subscribed")
)

// SubscriptionType represents an enumaration type for subscription type
type SubscriptionType string

// Values for SubscriptionType
const (
	Value    = "v"
	Movement = "m"
)

// Subscription implements un/subscription logic for crypobot
type Subscription struct {
	UserID          string
	Name            string
	Type            SubscriptionType
	AgainstCurrency string
	Account         *cryptobot.Account
}

// Subscribe subscribes for value change or account movement events for the given account
func subscribe(userID string, stype SubscriptionType, currencyID string, addrDesc string, currencyAgainst string) error {
	s := &Subscription{}

	switch stype {
	case Value:
	case Movement:
		s.Type = stype
		break
	default:
		return ErrSubscriptionType
	}

	subsKey := fmt.Sprintf("%s:%s:%s", userID, stype, addrDesc)

	if _, ok := subscriptions[subsKey]; ok {
		return ErrAlreadyExists
	}

	//TODO: use a factory function to create account
	a, err := cryptobot.NewAccount(currencyID, addrDesc, nil)
	if err != nil {
		return err
	}

	s.UserID = userID
	s.Account = a
	s.AgainstCurrency = currencyAgainst

	subscriptions[subsKey] = s

	return nil
}

// Unsubscribe unsubscribes for value change or account movement events for the given account
func unsubscribe(userID string, stype SubscriptionType, currencyID string, addrDesc string) error {
	switch stype {
	case Value:
	case Movement:
		break
	default:
		return ErrSubscriptionType
	}

	subsKey := fmt.Sprintf("%s:%s:%s", userID, stype, addrDesc)

	if _, ok := subscriptions[subsKey]; !ok {
		return ErrNotExists
	}

	delete(subscriptions, subsKey)

	return nil
}
