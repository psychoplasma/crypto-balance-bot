package mongodb

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/google/uuid"
	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// CollectionName is the name of Subscription collection
const CollectionName = "Subscriptions"

// SubscriptionRepository is MongoDB implementation of SubscriptionRepository
type SubscriptionRepository struct {
	client       *mongo.Client
	databaseURI  string
	session      mongo.Session
	sessionMutex *sync.Mutex
	subs         *mongo.Collection
	txOpts       *options.TransactionOptions
}

// NewSubscriptionRepository creates a new instance of SubscriptionRepository
func NewSubscriptionRepository(uri string, databaseName string) (*SubscriptionRepository, error) {
	repo := &SubscriptionRepository{
		databaseURI: uri,
		txOpts: options.Transaction().
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())).
			SetReadConcern(readconcern.Snapshot()),
		sessionMutex: new(sync.Mutex),
	}

	if err := repo.connect(); err != nil {
		return nil, err
	}

	repo.subs = repo.client.
		Database(databaseName).
		Collection(CollectionName)

	return repo, nil
}

// Begin starts a new session for ACID operation
func (r *SubscriptionRepository) Begin() error {
	log.Printf("Begin transaction")
	r.checkConnection()
	r.sessionMutex.Lock()

	session, err := r.client.StartSession()
	if err != nil {
		r.sessionMutex.Unlock()
		return err
	}

	r.session = session

	return nil
}

// Fail rollbacks to the state before this change.
// Notice that due to the fact that every operation is wrapped in
// MongoDB Callback-API (https://docs.mongodb.com/manual/core/transactions-in-applications/#callback-api)
// Here, just close the session and release the locked mutex
func (r *SubscriptionRepository) Fail() {
	log.Printf("Rollback transaction")
	defer r.sessionMutex.Unlock()
	r.session.EndSession(context.Background())
}

// Success finalizes ACID operation
// Notice that due to the fact that every operation is wrapped in
// MongoDB Callback-API (https://docs.mongodb.com/manual/core/transactions-in-applications/#callback-api)
// Here, just close the session and release the locked mutex
func (r *SubscriptionRepository) Success() {
	log.Printf("Finalize transaction")
	defer r.sessionMutex.Unlock()
	r.session.EndSession(context.Background())
}

// NextIdentity returns the next available identity
func (r *SubscriptionRepository) NextIdentity(userID string) string {
	return userID + ":" + uuid.New().String()
}

// Size returns the total number of subscriptions persited in the repository
func (r *SubscriptionRepository) Size() int64 {
	c, err := r.subs.EstimatedDocumentCount(context.Background(), nil)
	if err != nil {
		return -1
	}

	return c
}

// Get returns the subscription for the given subscription id
func (r *SubscriptionRepository) Get(id string) (*domain.Subscription, error) {
	s, err := r.applyOperation(func() (interface{}, error) {
		return r.get(id)
	})
	if err != nil {
		return nil, err
	}

	return ToDomain(s.(*Subscription)), err
}

// GetAllForUser returns all subscriptions for the given user id
func (r *SubscriptionRepository) GetAllForUser(userID string) ([]*domain.Subscription, error) {
	subs, err := r.applyOperation(func() (interface{}, error) {
		return r.getByUserID(userID)
	})
	if err != nil {
		return nil, err
	}

	return ToDomainSlice(subs.([]*Subscription)), nil
}

// GetAllActivatedMovements returns all activated movement subscriptions
func (r *SubscriptionRepository) GetAllActivatedMovements() ([]*domain.Subscription, error) {
	subs, err := r.applyOperation(func() (interface{}, error) {
		return r.getByTypeAndActivation(domain.MovementSubscription, true)
	})
	if err != nil {
		return nil, err
	}

	return ToDomainSlice(subs.([]*Subscription)), nil
}

// GetAllActivatedValues returns all activated value subscriptions
func (r *SubscriptionRepository) GetAllActivatedValues() ([]*domain.Subscription, error) {
	subs, err := r.applyOperation(func() (interface{}, error) {
		return r.getByTypeAndActivation(domain.ValueSubscription, true)
	})
	if err != nil {
		return nil, err
	}

	return ToDomainSlice(subs.([]*Subscription)), nil
}

// Save persists/updates the given subscription
func (r *SubscriptionRepository) Save(s *domain.Subscription) error {
	_, err := r.applyOperation(func() (interface{}, error) {
		return nil, r.replaceOrInsert(FromDomain(s))
	})

	return err
}

// Remove removes the given subscription from the persistance
func (r *SubscriptionRepository) Remove(s *domain.Subscription) error {
	_, err := r.applyOperation(func() (interface{}, error) {
		return nil, r.delete(s.ID())
	})

	return err
}

func (r *SubscriptionRepository) connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(r.databaseURI))
	if err != nil {
		return err
	}

	r.client = client

	return nil
}

func (r *SubscriptionRepository) checkConnection() {
	ctx := context.Background()
	if err := r.client.Ping(ctx, readpref.Primary()); err != nil {
		// If the connection is dead, try to reconnect
		if err := r.client.Connect(ctx); err != nil {
			// Do not handle nor bubble up this error
			// It must be the concern of the user of this repo
			panic(err)
		}
	}
}

func (r *SubscriptionRepository) disconnect() error {
	return r.client.Disconnect(context.Background())
}

func (r *SubscriptionRepository) applyOperation(op func() (interface{}, error)) (interface{}, error) {
	callback := func(sctx mongo.SessionContext) (interface{}, error) {
		return op()
	}

	return r.session.WithTransaction(context.Background(), callback, r.txOpts)
}

func (r *SubscriptionRepository) get(id string) (*Subscription, error) {
	s := &Subscription{}
	query := bson.M{"_id": id}

	if err := r.subs.FindOne(context.Background(), query).Decode(s); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return s, nil
}

func (r *SubscriptionRepository) getByUserID(userID string) ([]*Subscription, error) {
	ctx := context.Background()
	query := bson.M{"user_id": userID}

	cursor, err := r.subs.Find(ctx, query)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	subs := make([]*Subscription, 0)
	if err = cursor.All(ctx, &subs); err != nil {
		return nil, err
	}

	return subs, nil
}

func (r *SubscriptionRepository) getByTypeAndActivation(stype domain.SubscriptionType, activated bool) ([]*Subscription, error) {
	ctx := context.Background()
	query := bson.M{"type": stype, "activated": activated}

	cursor, err := r.subs.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	subs := make([]*Subscription, 0)
	if err = cursor.All(ctx, &subs); err != nil {
		return nil, err
	}

	return subs, nil
}

func (r *SubscriptionRepository) replaceOrInsert(s *Subscription) error {
	opts := options.Replace().SetUpsert(true)
	query := bson.M{"_id": s.ID}
	update := s

	res, err := r.subs.ReplaceOne(context.Background(), query, update, opts)
	if err != nil {
		return err
	}

	if res.UpsertedID != nil && (res.UpsertedID).(string) != s.ID {
		return fmt.Errorf("failed to save the subscription (%s)", s.ID)
	}

	return nil
}

func (r *SubscriptionRepository) delete(id string) error {
	query := bson.M{"_id": id}
	res, err := r.subs.DeleteOne(context.Background(), query)
	if err != nil {
		return err
	}

	if res.DeletedCount != 1 {
		return fmt.Errorf("failed to delete the subscription (%s)", id)
	}

	return nil
}

// Subscription represents a document in MongoDB corresponding to domain.Subscription
type Subscription struct {
	ID              string `bson:"_id" json:"_id"`
	UserID          string `bson:"user_id" json:"user_id"`
	Type            string `bson:"type" json:"type"`
	Activated       bool   `bson:"activated" json:"activated"`
	Currency        string `bson:"currency" json:"currency"`
	AgainstCurrency string `bson:"against_currency" json:"against_currency"`
	Account         string `bson:"account" json:"account"`
	Balance         string `bson:"balance" json:"balance"`
	BlockHeight     int    `bson:"block_height" json:"block_height"`
}

// FromDomain converts domain.Subscription model to a MongoDB document representation
func FromDomain(s *domain.Subscription) *Subscription {
	if s == nil {
		return nil
	}

	return &Subscription{
		ID:              s.ID(),
		UserID:          s.UserID(),
		Type:            string(s.Type()),
		Activated:       s.IsActivated(),
		Currency:        s.Currency().Symbol,
		AgainstCurrency: s.AgainstCurrency().Symbol,
		Account:         s.Account(),
		Balance:         s.Balance().String(),
		BlockHeight:     s.BlockHeight(),
	}
}

// ToDomain converts MongoDB document representation of Subscription to domain model
func ToDomain(s *Subscription) *domain.Subscription {
	if s == nil {
		return nil
	}

	balace, ok := new(big.Int).SetString(s.Balance, 10)
	if !ok {
		panic(fmt.Errorf("%s is not a valid bignumber representation", s.Balance))
	}

	sub, _ := domain.DeepCopySubscription(
		s.ID,
		s.UserID,
		domain.SubscriptionType(s.Type),
		s.Activated,
		s.Account,
		services.CurrencyFactory[s.Currency],
		services.CurrencyFactory[s.AgainstCurrency],
		balace,
		s.BlockHeight,
	)
	return sub
}

// ToDomainSlice converts slice of MongoDB documents to slice of domain models
func ToDomainSlice(subs []*Subscription) []*domain.Subscription {
	if subs == nil {
		return nil
	}

	domainSlice := make([]*domain.Subscription, len(subs))
	for i, s := range subs {
		domainSlice[i] = ToDomain(s)
	}

	return domainSlice
}
