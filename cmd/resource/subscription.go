package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/services"
)

// Subscription represents domain.Subscription for resource
type Subscription struct {
	ID                  string `json:"id"`
	UserID              string `json:"user_id"`
	Type                string `json:"type"`
	Activated           bool   `json:"is_activated"`
	Currency            string `json:"currency"`
	AgainstCurrency     string `json:"against_currency,omitempty"`
	Account             string `json:"account"`
	BlockHeight         uint64 `json:"last_updated_block_height"`
	TotalReceived       string `json:"total_received"`
	TotalSpent          string `json:"total_spent"`
	StartingBlockHeight uint64 `json:"starting_block_height"`
}

func fromDomain(s *domain.Subscription) *Subscription {
	if s == nil {
		return nil
	}

	return &Subscription{
		ID:                  s.ID(),
		UserID:              s.UserID(),
		Type:                string(s.Type()),
		Activated:           s.IsActivated(),
		Currency:            s.Currency().Symbol,
		AgainstCurrency:     s.AgainstCurrency().Symbol,
		Account:             s.Account(),
		BlockHeight:         s.BlockHeight(),
		StartingBlockHeight: s.StartingBlockHeight(),
		TotalReceived:       s.TotalReceived().String(),
		TotalSpent:          s.TotalSpent().String(),
	}
}

func fromDomainSlice(subs []*domain.Subscription) []*Subscription {
	if subs == nil {
		return nil
	}

	domainSlice := make([]*Subscription, len(subs))
	for i, s := range subs {
		domainSlice[i] = fromDomain(s)
	}

	return domainSlice
}

// GetAvailableAssets returns available assets that can be subscribed for
func GetAvailableAssets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(services.CurrencyFactory); err != nil {
		http.Error(w, fmt.Sprintf("cannot encode data to json, %s", err.Error()), http.StatusInternalServerError)
	}
}

// GetSubscription returns subscription details
func GetSubscription(w http.ResponseWriter, r *http.Request) {
	subscriptionID := mux.Vars(r)["id"]
	if subscriptionID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
	}

	s, err := subsApp.GetSubscription(subscriptionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(fromDomain(s)); err != nil {
		http.Error(w, fmt.Sprintf("corrupted subscription data, %s", err.Error()), http.StatusInternalServerError)
	}
}

// GetSubscriptionsForUser returns all the subscription for the given user
func GetSubscriptionsForUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userID"]
	if userID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
	}

	subs, err := subsApp.GetSubscriptionsForUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(fromDomainSlice(subs)); err != nil {
		http.Error(w, fmt.Sprintf("corrupted subscription data, %s", err.Error()), http.StatusInternalServerError)
	}
}
