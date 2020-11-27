package cryptobot_test

import (
	"math/big"
	"testing"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

func TestApply(t *testing.T) {
	addr := "test-addr-1"
	mv1 := &domain.AccountMovement{
		Address: addr,
		Changes: map[int][]*domain.BalanceChange{
			10: {
				{
					Amount: big.NewInt(5),
				},
			},
		},
	}

	a := domain.NewAccount(addr)
	initBalance := new(big.Int).Set(a.Balance())

	a.Apply(mv1)

	diff := big.NewInt(0).Sub(a.Balance(), initBalance)
	if diff.Cmp(big.NewInt(5)) != 0 {
		t.Fatalf("expected balance diff is %d but got %d", 5, diff.Int64())
	}
}

func TestApply_WithAlreadyAppliedMovements(t *testing.T) {
	addr := "test-addr-1"
	mv1 := &domain.AccountMovement{
		Address: addr,
		Changes: map[int][]*domain.BalanceChange{
			10: {
				{
					Amount: big.NewInt(5),
				},
			},
		},
	}

	mv2 := &domain.AccountMovement{
		Address: addr,
		Changes: map[int][]*domain.BalanceChange{
			10: {
				{
					Amount: big.NewInt(9),
				},
			},
		},
	}

	a := domain.NewAccount(addr)
	initBalance := new(big.Int).Set(a.Balance())

	a.Apply(mv1)
	a.Apply(mv2)

	diff := big.NewInt(0).Sub(a.Balance(), initBalance)
	if diff.Cmp(big.NewInt(5)) != 0 {
		t.Fatalf("expected balance diff is %d but got %d", 5, diff.Int64())
	}
}
