// +integration

package blockchaindotcom_test

import (
	"testing"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchain/blockchaindotcom"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/publisher/telegram"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/services"
)

func TestGetAccountMovements(t *testing.T) {
	blockNum := 183579
	api := blockchaindotcom.NewBitcoinAPI(blockchaindotcom.BitcoinTranslator{})

	mv, err := api.GetAccountMovements("1AJbsFZ64EpEfS5UAjAfcUG8pH8Jn3rn1F", blockNum)
	if err != nil {
		t.Fatal(err)
	}

	event := domain.NewAccountAssetsMovedEvent("subs_id", services.BTC, mv.Sort())
	t.Log(telegram.MovementFormatter(event))

	changesExistForBlock := false
	for blockHeight := range mv.Changes {
		if blockHeight >= blockNum {
			changesExistForBlock = true
			break
		}
	}

	if !changesExistForBlock {
		t.Fatalf("expected to have changes in block#%d but got nothing", blockNum)
	}
}
