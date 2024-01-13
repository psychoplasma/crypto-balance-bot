package main

import (
	"fmt"
	"log"
	"reflect"
	"time"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/application"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/concurrency"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/services"
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

const observeInterval = time.Second * 20
const exitTimeout = time.Second * 30
const maxParallelism = 1000

// ObserverOptions represents configurables for MovementObserver
type ObserverOptions struct {
	// BlockHeightMargin will cause fetching subscriptions where the difference between
	// the latest block height of the corresponding blockchain and the subscrition's
	// last updated block height is more than this margin
	// subscription.blockHeight < latestBlockHeight - margin
	BlockHeightMargin uint64
	// Sleep time inbetween every observal
	ObserveInterval time.Duration
	// Maximum number of goroutines for one observal
	MaxParallelism int
	// Timeout when stopping the observer
	ExitTimeout time.Duration
}

// MovementObserver observes for account movements
type MovementObserver struct {
	sa                *application.SubscriptionApplication
	w                 *concurrency.Worker
	p                 Publisher
	isObserving       bool
	observeInterval   time.Duration
	maxParallelism    int
	exitTimeout       time.Duration
	blockHeightMargin uint64
	currency          string
	cs                domain.CurrencyService
}

// NewMovementObserver creates a new instance of MovementObserver. It will panic if no service can be found for the given currency
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
		// If no BlockHeightMargin is provided, it will set to 0 by default
		o.blockHeightMargin = opt.BlockHeightMargin

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

	cs, ok := services.CurrencyServiceFactory[o.currency]
	if !ok {
		panic(fmt.Errorf("no service found for currency %s", o.currency))
	}
	o.cs = cs

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

	bh, err := o.cs.GetLatestBlockHeight()
	if err != nil {
		return err
	}

	subs, err := o.sa.GetSubscriptionsForCurrency(o.currency, bh-o.blockHeightMargin)
	if err != nil {
		return err
	}

	// And then check whether or not there is
	// a change in movement for each in parallel
	for _, s := range subs {
		// Create a local copy of the slice's items to pass to the worker's runner function
		cs := *s
		if _, err := o.w.Run(func() {
			if e := o.sa.CheckAndApplyAccountMovements(&cs); e != nil {
				log.Printf("error while observing: %s", e.Error())
			}
		}); err != nil {
			return err
		}
	}

	o.w.WaitAll()

	return nil
}
