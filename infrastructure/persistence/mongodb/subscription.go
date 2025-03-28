package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/google/uuid"
	domain "github.com/psychoplasma/crypto-balance-bot"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// CollectionName is the name of Subscription collection
const CollectionName = "Subscription"

// DocumentLimitsPerQuery limits query result to a certain number of documents
const DocumentLimitsPerQuery = 1000

// SubscriptionRepository is MongoDB implementation of SubscriptionRepository
type SubscriptionRepository struct {
	client       *mongo.Client
	session      mongo.Session
	sessionMutex *sync.Mutex
	subs         *mongo.Collection
	txOpts       *options.TransactionOptions
}

// NewSubscriptionRepository creates a new instance of SubscriptionRepository
func NewSubscriptionRepository() *SubscriptionRepository {
	repo := &SubscriptionRepository{
		txOpts: options.Transaction().
			SetWriteConcern(writeconcern.New(writeconcern.WMajority())).
			SetReadConcern(readconcern.Snapshot()),
		sessionMutex: new(sync.Mutex),
	}

	return repo
}

// Connect creates a connection to the given mongodb instance and the database
func (r *SubscriptionRepository) Connect(uri string, databaseName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	r.client = client
	r.subs = r.client.
		Database(databaseName).
		Collection(CollectionName)

	return nil
}

// Disconnect closes connection with the connected mongodb instance
func (r *SubscriptionRepository) Disconnect() error {
	return r.client.Disconnect(context.Background())
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

// GetAllForCurrency returns all subscriptions for the given currency
func (r *SubscriptionRepository) GetAllForCurrency(currencySymbol string, updatedBefore uint64) ([]*domain.Subscription, error) {
	subs, err := r.applyOperation(func() (interface{}, error) {
		return r.getByCurrency(currencySymbol, updatedBefore)
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
	query := bson.M{"userId": userID}

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

func (r *SubscriptionRepository) getByCurrency(symbol string, bh uint64) ([]*Subscription, error) {
	ctx := context.Background()
	opts := options.Find()
	opts.SetLimit(DocumentLimitsPerQuery)
	query := bson.M{
		"currency":    symbol,
		"blockHeight": bson.M{"$lt": bh},
	}

	cursor, err := r.subs.Find(ctx, query, opts)
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
	query := bson.M{"_id": s.ID}

	log.Println("Checking if record exists...")
	log.Printf("Subscription: ID=%s, UserID=%s, Currency=%s, Account=%s",
		s.ID, s.UserID, s.Currency, s.Account)

	err := r.subs.FindOne(context.Background(), query).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return r.insert(s)
		}
		return err
	}

	return r.replace(s)
}

func (r *SubscriptionRepository) replace(s *Subscription) error {
	query := bson.M{
		"userId":   s.UserID,
		"currency": s.Currency,
		"account":  s.Account,
	}
	// update := bson.M{
	// 	"blockHeight":         s.BlockHeight,
	// 	"totalReceived":       s.TotalReceived,
	// 	"totalSpent":          s.TotalSpent,
	// 	"startingBlockHeight": s.StartingBlockHeight,
	// 	"filters":             s.Filters,
	// }

	update := bson.D{
		{
			Key: "$set",
			Value: bson.M{
				"blockHeight":         s.BlockHeight,
				"totalReceived":       s.TotalReceived,
				"totalSpent":          s.TotalSpent,
				"startingBlockHeight": s.StartingBlockHeight,
				"filters":             s.Filters,
			},
		},
	}

	log.Println("Updating record...")

	res, err := r.subs.ReplaceOne(context.Background(), query, update)
	if err != nil {
		return err
	}

	if res.MatchedCount < 1 || (res.UpsertedID).(string) != s.ID {
		return fmt.Errorf("failed to save the subscription (%s)", s.ID)
	}

	return nil
}

func (r *SubscriptionRepository) insert(s *Subscription) error {
	log.Println("Creating new record...")

	res, err := r.subs.InsertOne(context.Background(), s)
	if err != nil {
		return err
	}

	if res.InsertedID != nil && (res.InsertedID).(string) != s.ID {
		return fmt.Errorf("failed to create subscription (%s)", s.ID)
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
	ID                  string   `bson:"_id"                 json:"_id"`
	UserID              string   `bson:"userId"              json:"userId"`
	Currency            string   `bson:"currency"            json:"currency"`
	CurrencyDecimal     string   `bson:"currencyDecimal"     json:"currencyDecimal"`
	Account             string   `bson:"account"             json:"account"`
	BlockHeight         uint64   `bson:"blockHeight"         json:"blockHeight"`
	TotalReceived       string   `bson:"totalReceived"       json:"totalReceived"`
	TotalSpent          string   `bson:"totalSpent"          json:"totalSpent"`
	StartingBlockHeight uint64   `bson:"startingBlockHeight" json:"startingBlockHeight"`
	Filters             []Filter `bson:"filters"             json:"filters"`
}

// Filter represents a document in MongoDB corresponding to domain.Filter
type Filter struct {
	Condition string `bson:"condition" json:"condition"`
	IsMust    bool   `bson:"isMust"    json:"isMust"`
	Type      string `bson:"type"      json:"type"`
}

// FromDomain converts domain.Subscription model to a MongoDB document representation
func FromDomain(s *domain.Subscription) *Subscription {
	if s == nil {
		return nil
	}

	filters := []Filter{}
	for _, f := range s.Filters() {
		data, err := f.SerializeCondition()
		if err != nil {
			panic(err)
		}

		filters = append(filters, Filter{
			Condition: string(data),
			Type:      string(f.Type()),
			IsMust:    f.IsMust(),
		})
	}

	return &Subscription{
		ID:                  s.ID(),
		UserID:              s.UserID(),
		Currency:            s.Currency().Symbol,
		CurrencyDecimal:     s.Currency().Decimal.String(),
		Account:             s.Account(),
		BlockHeight:         s.BlockHeight(),
		Filters:             filters,
		StartingBlockHeight: s.StartingBlockHeight(),
		TotalReceived:       s.TotalReceived().String(),
		TotalSpent:          s.TotalSpent().String(),
	}
}

// ToDomain converts MongoDB document representation of Subscription to domain model
func ToDomain(s *Subscription) *domain.Subscription {
	if s == nil {
		return nil
	}

	totalReceived, ok := new(big.Int).SetString(s.TotalReceived, 10)
	if !ok {
		panic(fmt.Errorf("TotalReceived (%s) is not a valid bignumber representation", s.TotalReceived))
	}

	totalSpent, ok := new(big.Int).SetString(s.TotalSpent, 10)
	if !ok {
		panic(fmt.Errorf("TotalSpent (%s) is not a valid bignumber representation", s.TotalSpent))
	}

	decimal, ok := new(big.Int).SetString(s.CurrencyDecimal, 10)
	if !ok {
		panic(fmt.Errorf("CurrencyDecimal (%s) is not a valid bignumber representation", s.CurrencyDecimal))
	}

	filters := []*domain.Filter{}
	for _, f := range s.Filters {
		filter := domain.NewFilter(domain.FilterType(f.Type), nil, f.IsMust)
		err := filter.DeserializeCondition([]byte(f.Condition))
		if err != nil {
			panic(err)
		}

		filters = append(filters, filter)
	}

	sub, _ := domain.DeepCopySubscription(
		s.ID,
		s.UserID,
		s.Account,
		domain.Currency{
			Symbol:  s.Currency,
			Decimal: decimal,
		},
		filters,
		totalReceived,
		totalSpent,
		s.BlockHeight,
		s.StartingBlockHeight,
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
