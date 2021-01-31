// +integration

package blockbook_test

import (
	"math/big"
	"testing"

	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchain/blockbook"
)

const bitcoinHostURL = "https://btc1.trezor.io"
const ethereumHostURL = "https://eth1.trezor.io"

func TestGetAccountMovements_Bitcoin(t *testing.T) {
	address := "1AJbsFZ64EpEfS5UAjAfcUG8pH8Jn3rn1F"
	since := 183579
	expectedChanges := map[int]*big.Int{
		183579: big.NewInt(-4678300000),
		643714: big.NewInt(2413000),
	}

	api := blockbook.NewAPI(bitcoinHostURL, blockbook.BitcoinTranslator{})
	mv, err := api.GetAccountMovements(address, since)
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

func TestGetAccountMovements_Bitcoin_WithPages(t *testing.T) {
	// address := "1AJbsFZ64EpEfS5UAjAfcUG8pH8Jn3rn1F"

	address := "1NDyJtNTjmwk5xPNhjgAMu4HDHigtobu1s"
	since := 0
	expectedBalance := big.NewInt(2413000)
	pagingLimit := 1000

	api := blockbook.NewAPI(bitcoinHostURL, blockbook.BitcoinTranslator{}, &pagingLimit)
	mv, err := api.GetAccountMovements(address, since)
	if err != nil {
		t.Fatal(err)
	}

	balance := big.NewInt(0)
	for _, chs := range mv.Changes {
		blockChange := big.NewInt(0)
		for _, ch := range chs {
			blockChange.Add(blockChange, ch.Amount)
			balance.Add(balance, ch.Amount)
		}
	}

	if expectedBalance.Cmp(balance) != 0 {
		t.Fatalf("expected balance is %s but got %s", expectedBalance.String(), balance.String())
	}
}

func TestGetAccountMovements_Ethereum(t *testing.T) {
	address := "0x7EF5A6135f1FD6a02593eEdC869c6D41D934aef8"
	since := 8676237
	expectedChanges := map[int][]*big.Int{
		8676237: {
			big.NewInt(3691368),
			big.NewInt(-1476547215),
			big.NewInt(-0),
		},
		8676239: {big.NewInt(-3152535117)},
	}

	api := blockbook.NewAPI(ethereumHostURL, blockbook.EthereumTranslator{})
	mv, err := api.GetAccountMovements(address, since)
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

		if len(expectedChanges[blockHeight]) != len(expectedChanges[blockHeight]) {
			t.Fatalf("expected %d number of changes at block#%d but got %d",
				len(expectedChanges[blockHeight]), blockHeight, len(c))
		}
	}
}
