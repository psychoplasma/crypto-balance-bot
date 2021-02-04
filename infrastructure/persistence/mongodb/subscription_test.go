// +integration
package mongodb_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/persistence/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "TestIntegration"
const dbURI = "mongodb://127.0.0.1:27017"

var testSubs = []*domain.Subscription{}

func TestSubscriptionRepository_NextIdentity(t *testing.T) {
	cleanUp := helperCreateAndPopulateDB(t)
	defer cleanUp()

	r, err := mongodb.NewSubscriptionRepository(dbURI, dbName)
	if err != nil {
		t.Fatal(err)
	}

	if id := r.NextIdentity("user1"); id == "" {
		t.Fatalf("expected non-empty string, but got empty one")
	}
}

func TestSubscriptionRepository_Size(t *testing.T) {
	cleanUp := helperCreateAndPopulateDB(t)
	defer cleanUp()

	r, err := mongodb.NewSubscriptionRepository(dbURI, dbName)
	if err != nil {
		t.Fatal(err)
	}

	if r.Size() != int64(len(testSubs)) {
		t.Fatalf("expected size %d, but got %d", len(testSubs), r.Size())
	}
}

func TestSubscriptionRepository_Get(t *testing.T) {
	cleanUp := helperCreateAndPopulateDB(t)
	defer cleanUp()

	r, err := mongodb.NewSubscriptionRepository(dbURI, dbName)
	if err != nil {
		t.Fatal(err)
	}

	testItem := testSubs[0]

	if err := r.Begin(); err != nil {
		t.Fatal(err)
	}
	s, _ := r.Get(testItem.ID())
	r.Success()

	if s.ID() != testItem.ID() || s.UserID() != testItem.UserID() {
		t.Fatalf("expected (ID, UserID) (%s, %s), but got (%s, %s)",
			testItem.ID(), testItem.UserID(),
			s.ID(), s.UserID(),
		)
	}
}

func TestSubscriptionRepository_GetAllForUser(t *testing.T) {
	cleanUp := helperCreateAndPopulateDB(t)
	defer cleanUp()

	r, err := mongodb.NewSubscriptionRepository(dbURI, dbName)
	if err != nil {
		t.Fatal(err)
	}

	testUser := testSubs[0].UserID()
	expectedCount := 2

	if err := r.Begin(); err != nil {
		t.Fatal(err)
	}
	subs, _ := r.GetAllForUser(testUser)
	r.Success()

	if len(subs) != expectedCount {
		t.Fatalf("expected size %d, but got %d", expectedCount, len(subs))
	}
}

func TestSubscriptionRepository_GetAllForType(t *testing.T) {
	cleanUp := helperCreateAndPopulateDB(t)
	defer cleanUp()

	r, err := mongodb.NewSubscriptionRepository(dbURI, dbName)
	if err != nil {
		t.Fatal(err)
	}

	expectedCount := 3

	if err := r.Begin(); err != nil {
		t.Fatal(err)
	}
	subs, _ := r.GetAllForType(domain.MovementSubscription)
	r.Success()

	if len(subs) != expectedCount {
		t.Fatalf("expected size %d, but got %d", expectedCount, len(subs))
	}

	expectedCount = 2

	if err := r.Begin(); err != nil {
		t.Fatal(err)
	}
	subs, _ = r.GetAllForType(domain.ValueSubscription)
	r.Success()

	if len(subs) != expectedCount {
		t.Fatalf("expected size %d, but got %d", expectedCount, len(subs))
	}
}

func TestSubscriptionRepository_GetAllForCurrency(t *testing.T) {
	cleanUp := helperCreateAndPopulateDB(t)
	defer cleanUp()

	r, err := mongodb.NewSubscriptionRepository(dbURI, dbName)
	if err != nil {
		t.Fatal(err)
	}

	expectedCount := 3

	if err := r.Begin(); err != nil {
		t.Fatal(err)
	}
	subs, _ := r.GetAllForCurrency("eth")
	r.Success()

	if len(subs) != expectedCount {
		t.Fatalf("expected size %d, but got %d", expectedCount, len(subs))
	}

	expectedCount = 2

	if err := r.Begin(); err != nil {
		t.Fatal(err)
	}
	subs, _ = r.GetAllForCurrency("btc")
	r.Success()

	if len(subs) != expectedCount {
		t.Fatalf("expected size %d, but got %d", expectedCount, len(subs))
	}
}

func TestSubscriptionRepository_Save(t *testing.T) {
	cleanUp := helperCreateAndPopulateDB(t)
	defer cleanUp()

	r, err := mongodb.NewSubscriptionRepository(dbURI, dbName)
	if err != nil {
		t.Fatal(err)
	}

	expectedSize := r.Size() + 1
	testItem, _ := domain.NewSubscription(r.NextIdentity("user3"), "user3", domain.MovementSubscription, "account-6", domain.Currency{}, domain.Currency{}, 0)

	if err := r.Begin(); err != nil {
		t.Fatal(err)
	}
	r.Save(testItem)
	r.Success()

	if r.Size() != expectedSize {
		t.Fatalf("expected size %d, but got %d", expectedSize, r.Size())
	}

	s, _ := r.Get(testItem.ID())

	if s.ID() != testItem.ID() || s.UserID() != testItem.UserID() {
		t.Fatalf("expected (ID, UserID) (%s, %s), but got (%s, %s)",
			testItem.ID(), testItem.UserID(),
			s.ID(), s.UserID(),
		)
	}
}

func TestSubscriptionRepository_Remove(t *testing.T) {
	cleanUp := helperCreateAndPopulateDB(t)
	defer cleanUp()

	r, err := mongodb.NewSubscriptionRepository(dbURI, dbName)
	if err != nil {
		t.Fatal(err)
	}

	expectedSize := r.Size() - 1
	testItem := testSubs[0]

	if err := r.Begin(); err != nil {
		t.Fatal(err)
	}
	r.Remove(testItem)
	r.Success()

	if r.Size() != expectedSize {
		t.Fatalf("expected size %d, but got %d", expectedSize, r.Size())
	}

	s, _ := r.Get(testItem.ID())

	if s != nil {
		t.Fatalf("expected subscription item nil, but got %#v", s)
	}
}

func helperReadTestData(t *testing.T, filename string) []*mongodb.Subscription {
	f, err := os.Open(filepath.Join("./testdata", filename))
	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	td := []*mongodb.Subscription{}
	if err := json.Unmarshal(data, &td); err != nil {
		t.Fatal(err)
	}
	testSubs = mongodb.ToDomainSlice(td)

	return td
}

func helperCreateAndPopulateDB(t *testing.T) func() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		t.Fatal(err)
	}

	db := client.Database(dbName)
	if db == nil {
		t.Fatal(fmt.Errorf("%s database doesn't exist", dbName))
	}

	coll := db.Collection(mongodb.CollectionName)
	if coll != nil {
		// t.Fatal(fmt.Errorf("%s collection doesn't exist", mongodb.CollectionName))
		if err := db.Drop(context.Background()); err != nil {
			t.Fatal(err)
		}

		if err := db.CreateCollection(ctx, mongodb.CollectionName); err != nil {
			t.Fatal(err)
		}
	}

	// Create the corresponding table if it does not exist in test db
	td := helperReadTestData(t, "subscriptions.json")
	for _, d := range td {
		res, err := coll.InsertOne(ctx, d)
		if err != nil {
			t.Fatal(err)
		}
		if res.InsertedID != d.ID {
			t.Fatal(fmt.Errorf("inserted document is corrupted"))
		}
	}

	// Return a function for clean-up
	dbCleanUp := func() {
		if err := db.Drop(context.Background()); err != nil {
			t.Fatal(err)
		}
	}

	return dbCleanUp
}
