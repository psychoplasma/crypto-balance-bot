package cryptobot_test

import (
	"math/big"
	"reflect"
	"testing"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/services"
)

func TestApply(t *testing.T) {
	addr := "test-addr-1"
	mv1 := domain.NewAccountMovements(addr)
	mv1.Receive(10, 1613721092, "txhash-test1", big.NewInt(5), "addr-sender")

	s, err := domain.NewSubscription("sub-1", "user-1", addr, services.ETH, 0)
	if err != nil {
		t.Fatal(err)
	}
	initReceivedBalance := new(big.Int).Set(s.TotalReceived())
	subscriber := NewMockEventSubscriber(
		reflect.TypeOf(new(domain.AccountAssetsMovedEvent)))
	domain.DomainEventPublisherInstance().Reset()
	domain.DomainEventPublisherInstance().Subscribe(subscriber)

	s.ApplyMovements(mv1)

	diff := new(big.Int).Sub(s.TotalReceived(), initReceivedBalance)
	if diff.Cmp(big.NewInt(5)) != 0 {
		t.Fatalf("expected balance diff is %d but got %d", 5, diff.Int64())
	}

	if !subscriber.IsEventHandled() {
		t.Fatal("expected to publish an AccountAssetsMovedEvent but got nothing")
	}
}

func TestApply_WithAlreadyAppliedMovements(t *testing.T) {
	addr := "test-addr-1"
	mv1 := domain.NewAccountMovements(addr)
	mv1.Receive(10, 1613721092, "txhash-test1", big.NewInt(5), "addr-sender")

	s, err := domain.NewSubscription("sub-1", "user-1", addr, services.ETH, 0)
	if err != nil {
		t.Fatal(err)
	}
	initReceivedBalance := new(big.Int).Set(s.TotalReceived())
	eventSubs := NewMockEventSubscriber(
		reflect.TypeOf(new(domain.AccountAssetsMovedEvent)))
	domain.DomainEventPublisherInstance().Reset()
	domain.DomainEventPublisherInstance().Subscribe(eventSubs)

	s.ApplyMovements(mv1)
	if !eventSubs.IsEventHandled() {
		t.Fatal("expected to publish an AccountAssetsMovedEvent but got nothing")
	}

	mv2 := domain.NewAccountMovements(addr)
	mv2.Receive(9, 1613721092, "txhash-test1", big.NewInt(9), "addr-sender")
	eventSubs.Reset()
	domain.DomainEventPublisherInstance().Reset()
	domain.DomainEventPublisherInstance().Subscribe(eventSubs)
	s.ApplyMovements(mv2)
	if eventSubs.IsEventHandled() {
		t.Fatal("expected not to publish any event but got an AccountAssetsMovedEvent")
	}

	diff := new(big.Int).Sub(s.TotalReceived(), initReceivedBalance)
	if diff.Cmp(big.NewInt(5)) != 0 {
		t.Fatalf("expected balance diff is %d but got %d", 5, diff.Int64())
	}
}

func TestFilter(t *testing.T) {

}
