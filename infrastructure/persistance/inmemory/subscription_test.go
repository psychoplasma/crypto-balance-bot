package inmemory_test

import (
	"os"
	"testing"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/persistance/inmemory"
)

var subsRepo = inmemory.NewSubscriptionReposititory()
var testSubs = []*domain.Subscription{}

func populateTestData() {
	s1, _ := domain.NewSubscription("1", "user1", "", domain.MovementSubscription, domain.Currency{})
	s1.Activate()

	s2, _ := domain.NewSubscription("2", "user1", "", domain.MovementSubscription, domain.Currency{})

	s3, _ := domain.NewSubscription("3", "user2", "", domain.MovementSubscription, domain.Currency{})
	s3.Activate()

	s4, _ := domain.NewSubscription("4", "user2", "", domain.MovementSubscription, domain.Currency{})
	s4.Activate()

	testSubs = append(testSubs, s1, s2, s3, s4)
	subsRepo.Add(s1)
	subsRepo.Add(s2)
	subsRepo.Add(s3)
	subsRepo.Add(s4)
}

func TestMain(m *testing.M) {
	populateTestData()
	os.Exit(m.Run())
}

func TestSize(t *testing.T) {
	if subsRepo.Size() != len(testSubs) {
		t.Fatalf("expected size %d, but got %d", len(testSubs), subsRepo.Size())
	}
}

func TestGet(t *testing.T) {
	testItem := testSubs[0]
	s, _ := subsRepo.Get(testItem.ID())

	if s.ID() != testItem.ID() || s.UserID() != testItem.UserID() {
		t.Fatalf("expected (ID, UserID) (%s, %s), but got (%s, %s)",
			testItem.ID(), testItem.UserID(),
			s.ID(), s.UserID(),
		)
	}
}

func TestGetAllForUser(t *testing.T) {
	testUser := testSubs[0].UserID()
	expectedCount := 2
	subs, _ := subsRepo.GetAllForUser(testUser)

	if len(subs) != expectedCount {
		t.Fatalf("expected size %d, but got %d", expectedCount, len(subs))
	}
}

func TestGetAllAcivated(t *testing.T) {
	expectedCount := 3
	subs, _ := subsRepo.GetAllActivated()

	if len(subs) != expectedCount {
		t.Fatalf("expected size %d, but got %d", expectedCount, len(subs))
	}
}

func TestAdd(t *testing.T) {
	expectedSize := len(testSubs) + 1
	testItem, _ := domain.NewSubscription("5", "user3", "", domain.MovementSubscription, domain.Currency{})

	subsRepo.Add(testItem)

	if subsRepo.Size() != expectedSize {
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

func TestAdd_ExistingSubscription_WithDifferentUserID(t *testing.T) {
	testItem, _ := domain.NewSubscription("1", "user3", "", domain.MovementSubscription, domain.Currency{})

	if err := subsRepo.Add(testItem); err == nil {
		t.Fatalf("expecting an error, but got nothing")
	}
}

func TestRemove(t *testing.T) {
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
