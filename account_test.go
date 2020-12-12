package cryptobot_test

import (
	"math/big"
	"reflect"
	"testing"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/services"
)

func TestApply(t *testing.T) {
	subsID := "test-subID-1"
	addr := "test-addr-1"
	mv1 := domain.NewAccountMovements(addr)
	mv1.AddBalanceChange(10, "txhash-test1", big.NewInt(5))

	a := domain.NewAccount(subsID, addr, services.ETH)
	initBalance := new(big.Int).Set(a.Balance())
	eventSubs := NewMockEventSubscriber(
		reflect.TypeOf(new(domain.AccountAssetsMovedEvent)))
	domain.DomainEventPublisherInstance().Reset()
	domain.DomainEventPublisherInstance().Subscribe(eventSubs)

	a.Apply(mv1)

	diff := new(big.Int).Sub(a.Balance(), initBalance)
	if diff.Cmp(big.NewInt(5)) != 0 {
		t.Fatalf("expected balance diff is %d but got %d", 5, diff.Int64())
	}

	if !eventSubs.IsEventHandled() {
		t.Fatal("expected to publish an AccountAssetsMovedEvent but got nothing")
	}
}

func TestApply_WithAlreadyAppliedMovements(t *testing.T) {
	subsID := "test-subID-1"
	addr := "test-addr-1"

	mv1 := domain.NewAccountMovements(addr)
	mv1.AddBalanceChange(10, "txhash-test1", big.NewInt(5))
	a := domain.NewAccount(subsID, addr, services.ETH)
	initBalance := new(big.Int).Set(a.Balance())
	eventSubs := NewMockEventSubscriber(
		reflect.TypeOf(new(domain.AccountAssetsMovedEvent)))
	domain.DomainEventPublisherInstance().Reset()
	domain.DomainEventPublisherInstance().Subscribe(eventSubs)

	a.Apply(mv1)
	if !eventSubs.IsEventHandled() {
		t.Fatal("expected to publish an AccountAssetsMovedEvent but got nothing")
	}

	mv2 := domain.NewAccountMovements(addr)
	mv2.AddBalanceChange(10, "txhash-test1", big.NewInt(9))
	eventSubs = NewMockEventSubscriber(
		reflect.TypeOf(new(domain.AccountAssetsMovedEvent)))
	domain.DomainEventPublisherInstance().Reset()
	domain.DomainEventPublisherInstance().Subscribe(eventSubs)
	a.Apply(mv2)
	if eventSubs.IsEventHandled() {
		t.Fatal("expected not to publish any event but got an AccountAssetsMovedEvent")
	}

	diff := new(big.Int).Sub(a.Balance(), initBalance)
	if diff.Cmp(big.NewInt(5)) != 0 {
		t.Fatalf("expected balance diff is %d but got %d", 5, diff.Int64())
	}
}

type MockEventSubscriber struct {
	eventHandled bool
	eventType    reflect.Type
}

func NewMockEventSubscriber(eventType reflect.Type) *MockEventSubscriber {
	return &MockEventSubscriber{
		eventHandled: false,
		eventType:    eventType,
	}
}

func (mes *MockEventSubscriber) HandleEvent(e interface{}) {
	mes.eventHandled = true
}

func (mes *MockEventSubscriber) SubscribedToEventType() reflect.Type {
	return mes.eventType
}

func (mes *MockEventSubscriber) IsEventHandled() bool {
	return mes.eventHandled
}
