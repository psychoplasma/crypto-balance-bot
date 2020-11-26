package cryptobot_test

import (
	"math/big"
	"testing"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

func TestApply(t *testing.T) {
	mv1 := &domain.AccountMovement{
		BlockHeight: 10,
		Changes: []*domain.BalanceChange{
			{
				Amount: big.NewInt(5),
			},
		},
	}
	addr := "test-addr-1"
	c := domain.Currency{
		Symbol:  "btc",
		Decimal: big.NewInt(8),
	}
	a := domain.NewAccount(c, addr)
	initBalance := big.NewInt(a.Balance().Int64())

	a.Apply(mv1)

	diff := big.NewInt(0).Sub(a.Balance(), initBalance)
	if diff.Cmp(big.NewInt(5)) != 0 {
		t.Fatalf("expected balance diff is %d but got %d", 5, diff.Int64())
	}
}

func TestApply_WithAlreadyAppliedMovement(t *testing.T) {
	mv1 := &domain.AccountMovement{
		BlockHeight: 10,
		Changes: []*domain.BalanceChange{
			{
				Amount: big.NewInt(5),
			},
		},
	}
	mv2 := &domain.AccountMovement{
		BlockHeight: 10,
		Changes: []*domain.BalanceChange{
			{
				Amount: big.NewInt(5),
			},
		},
	}
	addr := "test-addr-1"
	c := domain.Currency{
		Symbol:  "btc",
		Decimal: big.NewInt(8),
	}
	a := domain.NewAccount(c, addr)
	initBalance := big.NewInt(a.Balance().Int64())

	a.Apply(mv1)
	a.Apply(mv2)

	diff := big.NewInt(0).Sub(a.Balance(), initBalance)
	if diff.Cmp(big.NewInt(5)) != 0 {
		t.Fatalf("expected balance diff is %d but got %d", 5, diff.Int64())
	}
}
