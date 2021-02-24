package telegram_test

import (
	"math/big"
	"testing"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/publisher/telegram"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/services"
)

func TestMovementFormatter(t *testing.T) {
	expectedString := "```\ntest1 Received\n{\n\taddr-sender\n\t0.005000 eth\n\ttime@2021-02-18T16:51:32+09:00\n\tblock#12\n}\ntest1 Spent\n{\n\taddr-receiver\n\t0.002000 eth\n\ttime@2021-02-18T16:51:32+09:00\n\tblock#12\n}\ntest1 Received\n{\n\taddr-sender\n\t0.009000 eth\n\ttime@2021-02-19T16:51:32+09:00\n\tblock#23\n}\n```"

	acms := domain.NewAccountMovements("test1")
	acms.Receive(12, 1613634692, "tx-hash-1", big.NewInt(5000000000000000), "addr-sender")
	acms.Spend(12, 1613634692, "tx-hash-1", big.NewInt(2000000000000000), "addr-receiver")
	acms.Receive(23, 1613721092, "tx-hash-3", big.NewInt(9000000000000000), "addr-sender")
	event := domain.NewAccountAssetsMovedEvent("test-subsID-1", acms.Address, services.ETH, acms.Transfers)

	s := telegram.MovementFormatter(event)

	t.Log(s)

	if s != expectedString {
		t.Fatalf("expected string is\n%s\nbut got\n%s", expectedString, s)
	}
}
