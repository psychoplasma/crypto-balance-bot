package application

import (
	"log"
	"time"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/concurrency"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/services"
)

// Publisher defines the functionalities for publishing services
type Publisher interface {
	PublishMessage(userID string, msg interface{})
}

// MovementObserver observes for account movements
type MovementObserver struct {
	sr domain.SubscriptionRepository
	w  *concurrency.Worker
	ps []Publisher
}

// NewMovementObserver creates a new instance of MovementObserver
func NewMovementObserver(sr domain.SubscriptionRepository) *MovementObserver {
	return &MovementObserver{
		sr: sr,
		ps: []Publisher{},
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

		time.Sleep(time.Second * 10)
	}
}

// RegisterPublisher registers a publisher
func (o *MovementObserver) RegisterPublisher(p Publisher) {
	o.ps = append(o.ps, p)
}

func (o *MovementObserver) observe() error {
	subs, err := o.sr.GetAllActivatedMovements()
	if err != nil {
		return err
	}

	// And then check whether or not there is
	// a change in movement for each in parallel
	for _, s := range subs {
		log.Printf("subs: %+v\n", s)
		// FIXME: shouldn't pass by reference, will cause data corruption in parallel processing
		ss := *s

		if _, err := o.w.Run(func() {
			o.notify(
				ss.UserID(),
				o.checkForAccountMovements(&ss))
		}); err != nil {
			return err
		}

		time.Sleep(time.Millisecond * 1000)
	}

	o.w.WaitAll()

	return nil
}

func (o *MovementObserver) checkForAccountMovements(s *domain.Subscription) interface{} {
	sm := domain.NewSubscriptionMovements(s.ID(), s.Currency())
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

		sm.AddAccountMovements(acm)
	}

	return sm
}

func (o *MovementObserver) notify(userID string, i interface{}) {
	for _, p := range o.ps {
		if p == nil {
			continue
		}
		p.PublishMessage(userID, i)
	}
}

func (o *MovementObserver) clear() {
	o.ps = []Publisher{}
}
