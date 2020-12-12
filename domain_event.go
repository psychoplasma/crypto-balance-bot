package cryptobot

import (
	"reflect"
	"sync"
	"time"
)

// DomainEvent represents common functionalities of an event occuring on doamin
type DomainEvent interface {
	EventVersion() int
	OccurredOn() time.Time
}

// AllDomainEvents represents all kind of domain events. This is not
// a real domain event, it's used to subscribe to all domain events
type AllDomainEvents struct {
}

// EventVersion returns event version
func (de *AllDomainEvents) EventVersion() int {
	return 0
}

// OccurredOn returns the time at when the event occurs
func (de *AllDomainEvents) OccurredOn() time.Time {
	return time.Now()
}

// DomainEventSubscriber represents common functionalities of a subscriber for domain events
type DomainEventSubscriber interface {
	HandleEvent(domainEvent interface{})
	SubscribedToEventType() reflect.Type // Returns type of the subscribed DomainEvent
}

// DomainEventPublisher represents publisher for domain events
type DomainEventPublisher struct {
	m            *sync.Mutex
	isPublishing bool
	subscribers  []DomainEventSubscriber
}

func newDomainEventPublisher() *DomainEventPublisher {
	dep := &DomainEventPublisher{
		m: &sync.Mutex{},
	}
	dep.setPublishing(false)
	dep.ensureSubscribersList()

	return dep
}

var depInstance *DomainEventPublisher

// DomainEventPublisherInstance returns a singleton instance
func DomainEventPublisherInstance() *DomainEventPublisher {
	if depInstance == nil {
		depInstance = newDomainEventPublisher()
	}

	return depInstance
}

// Publish publishes the given domain event to its specific subscribers
// Subscriber can subscribe for a specific event or any kind of events
func (dep *DomainEventPublisher) Publish(domainEvent interface{}) {
	if dep.isPublishing || !dep.hasSubscribers() {
		return
	}

	dep.setPublishing(true)
	defer dep.setPublishing(false)

	eType := reflect.TypeOf(domainEvent)

	for _, s := range dep.subscribers {
		sType := s.SubscribedToEventType()
		if eType == sType || sType == reflect.TypeOf(new(AllDomainEvents)) {
			s.HandleEvent(domainEvent)
		}
	}
}

// PublishAll publishes all the given domain events
func (dep *DomainEventPublisher) PublishAll(events []interface{}) {
	for _, event := range events {
		dep.Publish(event)
	}
}

// Reset clears the subscriber list
func (dep *DomainEventPublisher) Reset() {
	if !dep.isPublishing {
		dep.subscribers = nil
	}
}

// Subscribe subscribes a subscriber for a specific domain event
func (dep *DomainEventPublisher) Subscribe(des DomainEventSubscriber) {
	if !dep.isPublishing {
		dep.ensureSubscribersList()
		dep.subscribers = append(dep.subscribers, des)
	}
}

func (dep *DomainEventPublisher) ensureSubscribersList() {
	if !dep.hasSubscribers() {
		dep.subscribers = make([]DomainEventSubscriber, 0)
	}
}

func (dep *DomainEventPublisher) setPublishing(flag bool) {
	dep.m.Lock()
	defer dep.m.Unlock()

	dep.isPublishing = flag
}

func (dep *DomainEventPublisher) hasSubscribers() bool {
	return dep.subscribers != nil
}
