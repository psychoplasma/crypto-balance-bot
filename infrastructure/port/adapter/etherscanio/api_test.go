// +integration

package etherscanio_test

import (
	"testing"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/notification"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/etherscanio"
)

func TestGetTxsOfAddress(t *testing.T) {
	blockNum := 7000000
	api := etherscanio.NewEthereumAPI(etherscanio.EthereumTranslator{})

	mv, err := api.GetTxsOfAddress("0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae", blockNum)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(notification.MovementFormatter(map[*domain.Account][]*domain.AccountMovement{
		domain.NewAccount(domain.Currency{}, "0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae"): mv,
	}))

	if mv[0] == nil || mv[len(mv)-1].BlockHeight < blockNum {
		t.Fatalf("expected to have changes in block#%d but got nothing", blockNum)
	}
}
