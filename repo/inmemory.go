package repo

type SubscriptionInMemoryReposititory struct {
}

func SubscriptionInMemoryReposititory() SubscriptionInMemoryReposititory {
	return SubscriptionInMemoryReposititory{}
}

func (r *SubscriptionInMemoryReposititory) Size() int {
	return 0
}

func (r *SubscriptionInMemoryReposititory) Get(id string) (*Subscription, error) {
	return nil, nil
}

func (r *SubscriptionInMemoryReposititory) GetAllForUser(userID string) ([]*Subscription, error) {
	return nil, nil
}

func (r *SubscriptionInMemoryReposititory) Add(s *Subscription) error {
	return nil
}

func (r *SubscriptionInMemoryReposititory) Remove(id string) error {
	return nil
}
