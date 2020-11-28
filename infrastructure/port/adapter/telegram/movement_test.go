package telegram_test

import (
	"math/big"
	"testing"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/telegram"
)

func TestMovementFormatter(t *testing.T) {
	expectedString := "```\nbtc[test1]\n{\n\tblock#2{ => -0.001000 btc => 0.009000 btc }\n\tblock#1{ => 0.005000 btc => -0.002000 btc }\n}\n```"
	mvs := []*domain.AccountMovement{
		{
			Address: "test1",
			Currency: domain.Currency{
				Symbol:  "btc",
				Decimal: big.NewInt(1000),
			},
			Changes: map[int][]*domain.BalanceChange{
				1: {
					{
						Amount: big.NewInt(5),
					},
					{
						Amount: big.NewInt(-2),
					},
				},
				2: {
					{
						Amount: big.NewInt(-1),
					},
					{
						Amount: big.NewInt(9),
					},
				},
			},
		},
	}

	s := telegram.MovementFormatter(mvs)

	t.Log(s)

	if s != expectedString {
		t.Fatalf("expected string is\n%s\nbut got\n%s", expectedString, s)
	}
}
