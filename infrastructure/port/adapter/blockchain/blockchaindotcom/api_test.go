// +integration

package blockchaindotcom_test

import (
	"testing"

	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchain/blockchaindotcom"
)

func TestGetAccountMovements(t *testing.T) {
	blockNum := uint64(183579)
	api := blockchaindotcom.NewAPI(blockchaindotcom.BitcoinTranslator{})

	mv, err := api.GetAccountMovements("1AJbsFZ64EpEfS5UAjAfcUG8pH8Jn3rn1F", blockNum)
	if err != nil {
		t.Fatal(err)
	}

	changesExistForBlock := false
	for _, t := range mv.Transfers {
		if t.BlockHeight >= blockNum {
			changesExistForBlock = true
			break
		}
	}

	if !changesExistForBlock {
		t.Fatalf("expected to have changes in block#%d but got nothing", blockNum)
	}
}

func TestGetLatestBlockHeight(t *testing.T) {
	api := blockchaindotcom.NewAPI(blockchaindotcom.BitcoinTranslator{})
	bh, err := api.GetLatestBlockHeight()
	if err != nil {
		t.Fatal(err)
	}

	if bh == 0 {
		t.Fatal("expected anything but 0")
	}
}
