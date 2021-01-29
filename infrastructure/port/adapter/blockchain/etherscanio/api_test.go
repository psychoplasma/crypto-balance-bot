// +integration

package etherscanio_test

import (
	"testing"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchain/etherscanio"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/publisher/telegram"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/services"
)

func TestGetTxsOfAddress(t *testing.T) {
	blockNum := 11000000
	api := etherscanio.NewEthereumAPI(etherscanio.EthereumTranslator{})

	mv, err := api.GetTxsOfAddress("0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae", blockNum)
	if err != nil {
		t.Fatal(err)
	}

	event := domain.NewAccountAssetsMovedEvent("subs_id", services.ETH, mv.Sort())
	t.Log(telegram.MovementFormatter(event))

	numOfChanges := 0
	for blockHeight, chs := range mv.Changes {
		numOfChanges += len(chs)

		if blockHeight < blockNum {
			t.Fatalf("expected to have blocks higher than %d but got a block#%d", blockNum, blockHeight)
		}
	}

	if numOfChanges == 0 {
		t.Fatalf("expected to have changes since block#%d but got nothing", blockNum)
	}
}
