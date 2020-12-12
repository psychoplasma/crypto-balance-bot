package application

import (
	"log"
	"reflect"
	"time"

	domain "github.com/psychoplasma/crypto-balance-bot"
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
func NewAccountAssetMovedEventSubscriber(p Publisher, userID string) *AccountAssetMovedEventSubscriber {
	return &AccountAssetMovedEventSubscriber{
		p:      p,
		userID: userID,
	}
}

// HandleEvent sends telegram message about account movement to this user
func (s *AccountAssetMovedEventSubscriber) HandleEvent(event interface{}) {
	acms, b := event.(*domain.AccountAssetsMovedEvent)
	if !b {
		log.Printf("unexpected event type, %+v\n", event)
		return
	}
	s.p.PublishMessage(s.userID, acms)
}

// SubscribedToEventType returns type of AccountAssetsMovedEvent to subscribe for it
func (s *AccountAssetMovedEventSubscriber) SubscribedToEventType() reflect.Type {
	return reflect.TypeOf(new(domain.AccountAssetsMovedEvent))
}

const observeInterval = 10

// MovementObserver observes for account movements
type MovementObserver struct {
	sr domain.SubscriptionRepository
	w  *concurrency.Worker
	p  Publisher
}

// NewMovementObserver creates a new instance of MovementObserver
func NewMovementObserver(sr domain.SubscriptionRepository, p Publisher) *MovementObserver {
	return &MovementObserver{
		sr: sr,
		p:  p,
		w:  concurrency.NewWorker(100, time.Second*30),
	}
}

// Observe starts observing for changes and blocks the current working thread
func (o *MovementObserver) Observe() {
	log.Printf("Starting the MovementObserver for movement changes")

	for true {
		if err := o.observe(); err != nil {
			log.Printf("exiting MovementObserver with an error: %s", err.Error())
			break
		}

		time.Sleep(time.Second * observeInterval)
	}
}

func (o *MovementObserver) observe() error {
	defer domain.DomainEventPublisherInstance().Reset()

	subs, err := o.sr.GetAllActivatedMovements()
	if err != nil {
		return err
	}

	// And then check whether or not there is
	// a change in movement for each in parallel
	for _, s := range subs {
		log.Printf("subs: %+v\n", s)

		domain.DomainEventPublisherInstance().
			Subscribe(NewAccountAssetMovedEventSubscriber(o.p, s.UserID()))

		// FIXME: shouldn't pass by reference, will cause data corruption in parallel processing
		ss := *s

		if _, err := o.w.Run(func() {
			o.checkForAccountMovements(&ss)
		}); err != nil {
			return err
		}

		time.Sleep(time.Millisecond * 1000)
	}

	o.w.WaitAll()

	return nil
}

func (o *MovementObserver) checkForAccountMovements(s *domain.Subscription) {
	for _, a := range s.Accounts() {
		acm, err := services.
			CurrencyServiceFactory[a.Currency().Symbol].
			GetTxsOfAddress(a.Address(), a.BlockHeight()+1)
		if err != nil {
			// FIXME: do not expose any details of the subscription
			log.Printf("failed to fetch movements for address(%s), %s", a.Address(), err)
			continue
		}

		// TODO: Doesn't look right place to apply changes on the domain object. Maybe need a domain service???
		a.Apply(acm.Sort())
	}
}
