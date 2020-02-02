package main

import (
	"log"
	"math/big"
)

type notifier func(recipient string, msg interface{})

func checkBalances(n notifier) {
	for _, s := range subscriptions {
		go checkBalance(s, n)
	}
}

func checkBalance(s *Subscription, n notifier) {
	movements := s.Account.UpdateBalances()
	if len(movements) < 1 {
		return
	}

	msg, err := formatMovements(movements)
	if err != nil {
		log.Printf("cannot format account movements for user %s", s.UserID)
		return
	}

	n(s.UserID, msg)
}

// TODO: implement this
func formatMovements(m map[string]*big.Int) (string, error) {
	return "", nil
}
