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

		time.Sleep(time.Second * 30)
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
			if c := o.checkForAccountMovements(&ss); c != nil {
				o.notify(ss.UserID(), c)
			}
		}); err != nil {
			return err
		}
	}

	o.w.WaitAll()

	return nil
}

func (o *MovementObserver) checkForAccountMovements(s *domain.Subscription) interface{} {
	changes := make(map[*domain.Account][]*domain.AccountMovement)
	for _, a := range s.Accounts() {
		movements, err := services.
			CurrencyServiceFactory[a.Currency().Symbol].
			GetTxsOfAddress(a.Address(), a.BlockHeight()+1)
		if err != nil {
			log.Printf("failed to fetch movements for address(%s), %s", a.Address(), err)
			continue
		}

		for _, mv := range movements {
			// TODO: Doesn't look right place to apply changes on the domain object. Maybe need a domain service???
			a.Apply(mv)
		}

		changes[a] = movements
	}

	return changes
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
