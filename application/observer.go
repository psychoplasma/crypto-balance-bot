package application

import (
	"log"
	"time"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/concurrency"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/services"
)

// Publisher is
type Publisher interface {
	PublishMessage(userID string, msg interface{})
}

// Observer observes for account movements
type Observer struct {
	sr domain.SubscriptionRepository
	w  *concurrency.Worker
	ps []Publisher
}

// NewObserver creates a new instance of Observer
func NewObserver(sr domain.SubscriptionRepository) *Observer {
	return &Observer{
		sr: sr,
		ps: []Publisher{},
		w:  concurrency.NewWorker(100, time.Second*30),
	}
}

// Observe starts observing for changes and blocks the current working thread
func (o *Observer) Observe() {
	log.Printf("Starting the observer for movement changes")

	for true {
		if err := o.observe(); err != nil {
			log.Printf("exiting observer with an error: %s", err.Error())
			break
		}

		time.Sleep(time.Second * 5)
	}
}

// RegisterPublisher registers a publisher
func (o *Observer) RegisterPublisher(p Publisher) {
	o.ps = append(o.ps, p)
}

func (o *Observer) observe() error {
	// Get all the activated subscriptions from the repository
	subs, err := o.sr.GetAllActivated()
	if err != nil {
		return err
	}

	// And then check whether or not there is
	// a change in movement for each in parallel
	for _, s := range subs {
		if _, err := o.w.Run(func() {
			if c := o.checkForAccountMovements(s); c != nil {
				o.notify(s.UserID(), c)
			}
		}); err != nil {
			return err
		}
	}

	o.w.WaitAll()

	return nil
}

func (o *Observer) checkForAccountMovements(s *domain.Subscription) interface{} {
	// FIXME: subscriptions fecthed from repo should be movement subscriptions, and shoudn't need to check the type here
	if s.Type() != domain.MovementSubscription {
		return nil
	}

	changes := make(map[*domain.Account][]*domain.AccountMovement)
	for _, a := range s.Accounts() {
		movements, err := services.
			CurrencyServiceFactory[a.Currency().Symbol].
			GetTxsOfAddress(a.Address(), a.BlockHeight())
		if err != nil {
			log.Printf("failed to fetch movements for address(%s), %s", a.Address(), err)
			continue
		}

		if movements == nil || len(movements) == 0 {
			continue
		}

		// TODO: Doesn't look right place to apply changes on the domain object. Maybe need a domain service???
		a.Apply(movements)

		changes[a] = movements
	}

	return changes
}

func (o *Observer) notify(userID string, i interface{}) {
	for _, p := range o.ps {
		if p == nil {
			continue
		}
		p.PublishMessage(userID, i)
	}
}

func (o *Observer) clear() {
	o.ps = []Publisher{}
}
