package inmemory_test

import (
	"os"
	"testing"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/persistence/inmemory"
)

var subsRepo = inmemory.NewSubscriptionRepository()
var testSubs = []*domain.Subscription{}

func populateSubsData() {
	s1, _ := domain.NewSubscription("1", "user1", "account-1", domain.Currency{Symbol: "c1"}, 0)
	s2, _ := domain.NewSubscription("2", "user1", "account-2", domain.Currency{Symbol: "c2"}, 0)
	s3, _ := domain.NewSubscription("3", "user2", "account-3", domain.Currency{Symbol: "c1"}, 0)
	s4, _ := domain.NewSubscription("4", "user2", "account-4", domain.Currency{Symbol: "c2"}, 0)
	s5, _ := domain.NewSubscription("5", "user3", "account-5", domain.Currency{Symbol: "c1"}, 0)

	testSubs = append(testSubs, s1, s2, s3, s4, s5)
	subsRepo.Save(s1)
	subsRepo.Save(s2)
	subsRepo.Save(s3)
	subsRepo.Save(s4)
	subsRepo.Save(s5)
}

func TestMain(m *testing.M) {
	populateSubsData()
	os.Exit(m.Run())
}

func TestSubscriptionRepository_NextIdentity(t *testing.T) {
	if id := subsRepo.NextIdentity("user-1"); id == "" {
		t.Fatalf("expected non-empty string, but got empty one")
	}
}

func TestSubscriptionRepository_Size(t *testing.T) {
	if subsRepo.Size() != int64(len(testSubs)) {
		t.Fatalf("expected size %d, but got %d", len(testSubs), subsRepo.Size())
	}
}

func TestSubscriptionRepository_Get(t *testing.T) {
	testItem := testSubs[0]
	s, _ := subsRepo.Get(testItem.ID())

	if s.ID() != testItem.ID() || s.UserID() != testItem.UserID() {
		t.Fatalf("expected (ID, UserID) (%s, %s), but got (%s, %s)",
			testItem.ID(), testItem.UserID(),
			s.ID(), s.UserID(),
		)
	}
}

func TestSubscriptionRepository_GetAllForUser(t *testing.T) {
	testUser := testSubs[0].UserID()
	expectedCount := 2
	subs, _ := subsRepo.GetAllForUser(testUser)

	if len(subs) != expectedCount {
		t.Fatalf("expected size %d, but got %d", expectedCount, len(subs))
	}
}

func TestSubscriptionRepository_GetAllForCurrency(t *testing.T) {
	expectedCount := 3
	subs, _ := subsRepo.GetAllForCurrency("c1")

	if len(subs) != expectedCount {
		t.Fatalf("expected size %d, but got %d", expectedCount, len(subs))
	}

	expectedCount = 2
	subs, _ = subsRepo.GetAllForCurrency("c2")

	if len(subs) != expectedCount {
		t.Fatalf("expected size %d, but got %d", expectedCount, len(subs))
	}
}

func TestSubscriptionRepository_Save(t *testing.T) {
	expectedSize := len(testSubs) + 1
	testItem, _ := domain.NewSubscription("6", "user3", "account-6", domain.Currency{}, 0)

	subsRepo.Save(testItem)

	if subsRepo.Size() != int64(expectedSize) {
		t.Fatalf("expected size %d, but got %d", expectedSize, subsRepo.Size())
	}

	s, _ := subsRepo.Get(testItem.ID())

	if s.ID() != testItem.ID() || s.UserID() != testItem.UserID() {
		t.Fatalf("expected (ID, UserID) (%s, %s), but got (%s, %s)",
			testItem.ID(), testItem.UserID(),
			s.ID(), s.UserID(),
		)
	}
}

func TestSubscriptionRepository_Remove(t *testing.T) {
	expectedSize := subsRepo.Size() - 1
	testItem := testSubs[0]

	subsRepo.Remove(testItem)

	if subsRepo.Size() != expectedSize {
		t.Fatalf("expected size %d, but got %d", expectedSize, subsRepo.Size())
	}

	s, _ := subsRepo.Get(testItem.ID())

	if s != nil {
		t.Fatalf("expected subscription item nil, but got %#v", s)
	}
}
