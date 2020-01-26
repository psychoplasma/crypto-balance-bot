package cryptobot

import (
	"errors"
)

// SubscriptionType represents an enumaration type for subscription type
type SubscriptionType string

// Values for SubscriptionType
const (
	Value    = "1"
	Movement = "2"
)

var ErrSubscriptionType = errors.New("wrong subscription type")

// Subscription implements un/subscription logic for crypobot
type Subscription struct {
	UserID          string
	Name            string
	Type            SubscriptionType
	AgainstCurrency string
	Account         *Account
}
