package main

import (
	"log"
	"reflect"
	"time"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/application"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/concurrency"
)

// Publisher defines the functionalities for publishing services
type Publisher interface {
	PublishMessage(userID string, msg interface{})
}

// AccountAssetMovedEventSubscriber implements domain.DomainEventSubscriber interface
type AccountAssetMovedEventSubscriber struct {
	p      Publisher
	userID string
}

// NewAccountAssetMovedEventSubscriber creates a new instance of subscriber for AccountAssetMoved event
func NewAccountAssetMovedEventSubscriber(p Publisher) *AccountAssetMovedEventSubscriber {
	return &AccountAssetMovedEventSubscriber{
		p: p,
	}
}

// HandleEvent sends telegram message about account movement to this user
func (s *AccountAssetMovedEventSubscriber) HandleEvent(event interface{}) {
	acms, b := event.(*domain.AccountAssetsMovedEvent)
	if !b {
		log.Printf("unexpected event type, %+v\n", event)
		return
	}
	s.p.PublishMessage(domain.UserIDFrom(acms.SubscriptionID()), acms)
}

// SubscribedToEventType returns type of AccountAssetsMovedEvent to subscribe for it
func (s *AccountAssetMovedEventSubscriber) SubscribedToEventType() reflect.Type {
	return reflect.TypeOf(new(domain.AccountAssetsMovedEvent))
}

const observeInterval = time.Second * 10
const exitTimeout = time.Second * 30
const maxParallelism = 1000

// ObserverOptions represents configurables for MovementObserver
type ObserverOptions struct {
	ObserveInterval time.Duration // Sleep time inbetween every observal
	MaxParallelism  int           // Maximum number of goroutines for one observal
	ExitTimeout     time.Duration // Timeout when stopping the observer
}

// MovementObserver observes for account movements
type MovementObserver struct {
	sa              *application.SubscriptionApplication
	w               *concurrency.Worker
	p               Publisher
	isObserving     bool
	observeInterval time.Duration
	maxParallelism  int
	exitTimeout     time.Duration
	currency        string
}

// NewMovementObserver creates a new instance of MovementObserver
func NewMovementObserver(sa *application.SubscriptionApplication, p Publisher, currency string, opts ...*ObserverOptions) *MovementObserver {
	o := &MovementObserver{
		currency:        currency,
		sa:              sa,
		p:               p,
		observeInterval: observeInterval,
		exitTimeout:     exitTimeout,
		maxParallelism:  maxParallelism,
	}

	for _, opt := range opts {
		if opt.ObserveInterval != 0 {
			o.observeInterval = opt.ObserveInterval
		}
		if opt.ExitTimeout != 0 {
			o.exitTimeout = opt.ExitTimeout
		}
		if opt.MaxParallelism != 0 {
			o.maxParallelism = opt.MaxParallelism
		}
	}

	o.w = concurrency.NewWorker(o.maxParallelism, o.exitTimeout)

	return o
}

// Start starts observing for changes and blocks the current working thread
func (o *MovementObserver) Start() {
	log.Printf("Starting MovementObserver")

	o.isObserving = true
	for o.isObserving {
		if err := o.observe(); err != nil {
			log.Printf("error while observing: %s", err.Error())
		}

		time.Sleep(o.observeInterval)
	}
}

// Stop stops observing gracefully
func (o *MovementObserver) Stop() {
	o.isObserving = false
	o.w.Stop()
}

func (o *MovementObserver) observe() error {
	domain.DomainEventPublisherInstance().
		Subscribe(NewAccountAssetMovedEventSubscriber(o.p))
	defer domain.DomainEventPublisherInstance().Reset()

	subs, err := o.sa.GetSubscriptionsForCurrency(o.currency)
	if err != nil {
		return err
	}

	// And then check whether or not there is
	// a change in movement for each in parallel
	for _, s := range subs {
		// Create a local copy of the slice's items to pass to the worker's runner function
		cs := *s
		if _, err := o.w.Run(func() { o.sa.CheckAndApplyAccountMovements(&cs) }); err != nil {
			return err
		}
	}

	o.w.WaitAll()

	return nil
}
