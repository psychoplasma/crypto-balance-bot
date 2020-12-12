package telegram_test

import (
	"math/big"
	"testing"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/telegram"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/services"
)

func TestMovementFormatter(t *testing.T) {
	expectedString := "```\nbtc[test1]\n{\n\tblock#12{ => 0.005000 btc => -0.002000 btc }\n\tblock#23{ => -0.001000 btc => 0.009000 btc }\n}\n```"

	acms := domain.NewAccountMovements("test1")
	acms.AddBalanceChange(12, "tx-hash-1", big.NewInt(500000))
	acms.AddBalanceChange(12, "tx-hash-1", big.NewInt(-200000))
	acms.AddBalanceChange(23, "tx-hash-2", big.NewInt(-100000))
	acms.AddBalanceChange(23, "tx-hash-3", big.NewInt(900000))
	event := domain.NewAccountAssetsMovedEvent("test-subsID-1", services.BTC, acms)

	s := telegram.MovementFormatter(event)

	t.Log(s)

	if s != expectedString {
		t.Fatalf("expected string is\n%s\nbut got\n%s", expectedString, s)
	}
}
