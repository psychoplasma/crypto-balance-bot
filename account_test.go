package cryptobot_test

import (
	"math/big"
	"testing"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/services"
)

func TestApply(t *testing.T) {
	addr := "test-addr-1"
	mv1 := domain.NewAccountMovements(addr)
	mv1.AddBalanceChange(10, "txhash-test1", big.NewInt(5))

	a := domain.NewAccount(addr, services.ETH)
	initBalance := new(big.Int).Set(a.Balance())

	a.Apply(mv1)

	diff := new(big.Int).Sub(a.Balance(), initBalance)
	if diff.Cmp(big.NewInt(5)) != 0 {
		t.Fatalf("expected balance diff is %d but got %d", 5, diff.Int64())
	}
}

func TestApply_WithAlreadyAppliedMovements(t *testing.T) {
	addr := "test-addr-1"
	mv1 := domain.NewAccountMovements(addr)
	mv1.AddBalanceChange(10, "txhash-test1", big.NewInt(5))

	mv2 := domain.NewAccountMovements(addr)
	mv2.AddBalanceChange(10, "txhash-test1", big.NewInt(9))

	a := domain.NewAccount(addr, services.ETH)
	initBalance := new(big.Int).Set(a.Balance())

	a.Apply(mv1)
	a.Apply(mv2)

	diff := new(big.Int).Sub(a.Balance(), initBalance)
	if diff.Cmp(big.NewInt(5)) != 0 {
		t.Fatalf("expected balance diff is %d but got %d", 5, diff.Int64())
	}
}
