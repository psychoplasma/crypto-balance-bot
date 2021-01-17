package cryptobot_test

import (
	"reflect"
	"testing"
	"time"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

type MockEventSubscriber struct {
	numOfHandledEvents int
	eventHandled       bool
	eventType          reflect.Type
}

func NewMockEventSubscriber(eventType reflect.Type) *MockEventSubscriber {
	return &MockEventSubscriber{
		eventHandled: false,
		eventType:    eventType,
	}
}

func (mes *MockEventSubscriber) HandleEvent(e interface{}) {
	mes.eventHandled = true
	mes.numOfHandledEvents++
}

func (mes *MockEventSubscriber) SubscribedToEventType() reflect.Type {
	return mes.eventType
}

func (mes *MockEventSubscriber) IsEventHandled() bool {
	return mes.eventHandled
}

func (mes *MockEventSubscriber) Reset() {
	mes.eventHandled = false
}

type MockDomainEvent struct {
}

func (mde *MockDomainEvent) EventVersion() int {
	return 1
}

func (mde *MockDomainEvent) OccurredOn() time.Time {
	return time.Now()
}

func MockDomainEventType() reflect.Type {
	return reflect.TypeOf(new(MockDomainEvent))
}

func TestPublish(t *testing.T) {
	subscriber1 := NewMockEventSubscriber(MockDomainEventType())
	subscriber2 := NewMockEventSubscriber(reflect.TypeOf(new(domain.AccountAssetsMovedEvent)))
	subscriber3 := NewMockEventSubscriber(reflect.TypeOf(new(domain.AllDomainEvents)))

	domain.DomainEventPublisherInstance().Reset()
	domain.DomainEventPublisherInstance().Subscribe(subscriber1)
	domain.DomainEventPublisherInstance().Subscribe(subscriber2)
	domain.DomainEventPublisherInstance().Subscribe(subscriber3)
	domain.DomainEventPublisherInstance().Publish(new(MockDomainEvent))

	if !subscriber1.IsEventHandled() {
		t.Fatal("expected MockDomainEventType to be handled but it wasn't")
	}

	if subscriber2.IsEventHandled() {
		t.Fatal("expected the event not to be handled but it was")
	}

	if !subscriber3.IsEventHandled() {
		t.Fatal("expected the event to be handled but it wasn't")
	}
}

func TestPublish_MultipleEvents(t *testing.T) {
	subscriber := NewMockEventSubscriber(MockDomainEventType())

	domain.DomainEventPublisherInstance().Reset()
	domain.DomainEventPublisherInstance().Subscribe(subscriber)
	domain.DomainEventPublisherInstance().Publish(new(MockDomainEvent))
	if !subscriber.IsEventHandled() {
		t.Fatal("expected MockDomainEventType to be handled but it wasn't")
	}
	subscriber.Reset()

	domain.DomainEventPublisherInstance().Publish(new(MockDomainEvent))
	if !subscriber.IsEventHandled() {
		t.Fatal("expected MockDomainEventType to be handled but it wasn't")
	}

	if subscriber.numOfHandledEvents != 2 {
		t.Fatalf("expected number of handled event 2 but got %d", subscriber.numOfHandledEvents)
	}
}
