// +integration

package etherscanio_test

import (
	"testing"

	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchain/etherscanio"
)

func TestGetAccountMovements(t *testing.T) {
	blockNum := uint64(11000000)
	api := etherscanio.NewAPI(etherscanio.EthereumTranslator{})

	mv, err := api.GetAccountMovements("0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae", blockNum)
	if err != nil {
		t.Fatal(err)
	}

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

func TestGetLatestBlockHeight(t *testing.T) {
	api := etherscanio.NewAPI(etherscanio.EthereumTranslator{})
	bh, err := api.GetLatestBlockHeight()
	if err != nil {
		t.Fatal(err)
	}

	if bh == 0 {
		t.Fatal("expected anything but 0")
	}
}
