package inmemory_test

import (
	"os"
	"testing"

	cyrptoBot "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/repo/inmemory"
)

var subsRepo = inmemory.NewSubscriptionReposititory()
var testSubs = []*cyrptoBot.Subscription{
	&cyrptoBot.Subscription{
		ID:        "1",
		UserID:    "user1",
		Activated: true,
	},
	&cyrptoBot.Subscription{
		ID:     "2",
		UserID: "user1",
	},
	&cyrptoBot.Subscription{
		ID:        "3",
		UserID:    "user2",
		Activated: true,
	},
	&cyrptoBot.Subscription{
		ID:        "4",
		UserID:    "user2",
		Activated: true,
	},
}

func TestMain(m *testing.M) {
	for _, s := range testSubs {
		subsRepo.Add(s)
	}

	os.Exit(m.Run())
}

func TestSize(t *testing.T) {
	if subsRepo.Size() != len(testSubs) {
		t.Fatalf("expected size %d, but got %d", len(testSubs), subsRepo.Size())
	}
}

func TestGet(t *testing.T) {
	testItem := testSubs[0]
	s, _ := subsRepo.Get(testItem.ID)

	if s.ID != testItem.ID || s.UserID != testItem.UserID {
		t.Fatalf("expected (ID, UserID) (%s, %s), but got (%s, %s)",
			testItem.ID, testItem.UserID,
			s.ID, s.UserID,
		)
	}
}

func TestGetAllForUser(t *testing.T) {
	testUser := testSubs[0].UserID
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
	testItem := &cyrptoBot.Subscription{
		ID:     "5",
		UserID: "user3",
	}

	subsRepo.Add(testItem)

	if subsRepo.Size() != expectedSize {
		t.Fatalf("expected size %d, but got %d", expectedSize, subsRepo.Size())
	}

	s, _ := subsRepo.Get(testItem.ID)

	if s.ID != testItem.ID || s.UserID != testItem.UserID {
		t.Fatalf("expected (ID, UserID) (%s, %s), but got (%s, %s)",
			testItem.ID, testItem.UserID,
			s.ID, s.UserID,
		)
	}
}

func TestAdd_ExistingSubscription_WithDifferentUserID(t *testing.T) {
	testItem := &cyrptoBot.Subscription{
		ID:     "1",
		UserID: "user3",
	}

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

	s, _ := subsRepo.Get(testItem.ID)

	if s != nil {
		t.Fatalf("expected subscription item nil, but got %#v", s)
	}
}
