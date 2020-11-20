package application

import (
	"log"
	"time"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/concurrency"
)

type Notifier func(recipient string, msg interface{})

type Observer struct {
	as domain.AccountService
	sr domain.SubscriptionRepository
	w  *concurrency.Worker
	ns map[string]Notifier
}

// NewObserver creates a new instance of Observer
func NewObserver(
	as domain.AccountService,
	sr domain.SubscriptionRepository) *Observer {
	return &Observer{
		as: as,
		sr: sr,
		ns: make(map[string]Notifier),
		w:  concurrency.NewWorker(100, time.Second*30),
	}
}

// Observe starts to observe for changes
func (o *Observer) Observe() error {
	for true {
		log.Panicln("Observing...")
		if err := o.observe(); err != nil {
			break
		}

		time.Sleep(time.Second * 5)
	}

	return nil
}

// RegisterNotifier registers a notifier. Notifiers with the same id will be replaced
func (o *Observer) RegisterNotifier(id string, n Notifier) {
	o.ns[id] = n
}

func (o *Observer) observe() error {
	// Get all the activated subscriptions from the repository
	subs, err := o.sr.GetAllActivated()
	if err != nil {
		return err
	}

	// And then check whether or not there is a change
	// in value/movement for each of them in parallel
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
	if s.Type() == domain.ValueSubscription {
		// TODO: implement
		return nil
	}

	changes := make(map[string][]*domain.AccountMovement)
	if s.Type() == domain.MovementSubscription {
		for _, a := range s.Accounts() {
			movements, err := o.as.FetchAccountMovements(a)
			if err != nil {
				log.Println(err.Error())
			}

			if movements == nil || len(movements) == 0 {
				continue
			}

			changes[a.Address()] = movements
		}

		return changes
	}

	return nil
}

func (o *Observer) notify(userID string, i interface{}) {
	if n, exist := o.ns[userID]; exist && n != nil {
		n(userID, i)
	}
}

func (o *Observer) clear() {
	o.ns = make(map[string]Notifier)
}
