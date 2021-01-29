// +integration

package blockbook_test

import (
	"math/big"
	"testing"

	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchain/bitcoin/blockbook"
)

const hostURL = "https://btc1.trezor.io"

var api = blockbook.NewBitcoinAPI(hostURL, blockbook.BitcoinTranslator{})

func TestGetTxsOfAddress(t *testing.T) {
	since := 183579
	expectedChanges := map[int]*big.Int{
		183579: big.NewInt(-4678300000),
		643714: big.NewInt(2413000),
	}

	mv, err := api.GetTxsOfAddress("1AJbsFZ64EpEfS5UAjAfcUG8pH8Jn3rn1F", since)
	if err != nil {
		t.Fatal(err)
	}

	if len(mv.Changes) != len(expectedChanges) {
		t.Fatalf("expected to have %d changes but got %d changes", len(expectedChanges), len(mv.Changes))
	}

	for blockHeight, c := range mv.Changes {
		if expectedChanges[blockHeight] == nil {
			t.Fatalf("expected a change at block#%d but got nothing", blockHeight)
		}

		if expectedChanges[blockHeight].Cmp(c[0].Amount) != 0 {
			t.Fatalf("expected a change %s at block#%d but got %s",
				expectedChanges[blockHeight].String(), blockHeight, c[0].Amount.String())
		}
	}
}
