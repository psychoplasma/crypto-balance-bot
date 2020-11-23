package application

import (
	"log"
	"time"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/concurrency"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/services"
)

type Publisher interface {
	PublishMessage(userID string, msg interface{})
}

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
	for true {
		log.Printf("Observing...")
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
			if c := o.checkForChange(s); c != nil {
				o.notify(s.UserID(), c)
			}
		}); err != nil {
			return err
		}
	}

	o.w.WaitAll()

	return nil
}

func (o *Observer) checkForChange(s *domain.Subscription) interface{} {
	changes := make(map[*domain.Account][]*domain.AccountMovement)

	if s.Type() == domain.MovementSubscription {
		for _, a := range s.Accounts() {
			movements, err := o.fetchAccountMovements(a)
			if err != nil {
				log.Printf("failed to fetch movements for address(%s), %s", a.Address(), err)
				continue
			}

			if movements == nil || len(movements) == 0 {
				continue
			}

			changes[a] = movements
		}

		return changes
	}

	return nil
}

func (o *Observer) fetchAccountMovements(a *domain.Account) ([]*domain.AccountMovement, error) {
	currencyService, err := services.CurrencyFactory(a.Currency())
	if err != nil {
		return nil, err
	}

	ms, err := currencyService.GetAddressTxs(a.Address(), a.BlockHeight())
	if err != nil {
		return nil, err
	}

	return ms, nil
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
