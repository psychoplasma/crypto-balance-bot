// +integration

package blockbook_test

import (
	"math/big"
	"testing"

	domain "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchain/blockbook"
)

const bitcoinHostURL = "https://btc1.trezor.io"
const ethereumHostURL = "https://eth1.trezor.io"

func TestGetAccountMovements_Bitcoin(t *testing.T) {
	address := "1AJbsFZ64EpEfS5UAjAfcUG8pH8Jn3rn1F"
	since := uint64(183579)
	expectedTransfers := map[uint64]*domain.Transfer{
		183579: {
			Type:        domain.Spent,
			BlockHeight: 183579,
			Amount:      big.NewInt(4678300000),
		},
		643714: {
			Type:        domain.Received,
			BlockHeight: 643714,
			Amount:      big.NewInt(2413000),
		},
	}

	api := blockbook.NewAPI(bitcoinHostURL, blockbook.BitcoinTranslator{})
	mv, err := api.GetAccountMovements(address, since)
	if err != nil {
		t.Fatal(err)
	}

	if len(mv.Transfers) != len(expectedTransfers) {
		t.Fatalf("expected to have %d changes but got %d changes", len(expectedTransfers), len(mv.Transfers))
	}

	for _, tr := range mv.Transfers {
		if expectedTransfers[tr.BlockHeight].Value().Cmp(tr.Value()) != 0 {
			t.Fatalf("expected a change %s at block#%d but got %s",
				expectedTransfers[tr.BlockHeight].Value().String(), tr.BlockHeight, tr.Value().String())
		}
	}
}

func TestGetAccountMovements_Bitcoin_WithPages(t *testing.T) {
	address := "1AJbsFZ64EpEfS5UAjAfcUG8pH8Jn3rn1F"
	since := uint64(0)
	expectedBalance := big.NewInt(2413000)
	pagingLimit := 1000

	api := blockbook.NewAPI(bitcoinHostURL, blockbook.BitcoinTranslator{}, &pagingLimit)
	mv, err := api.GetAccountMovements(address, since)
	if err != nil {
		t.Fatal(err)
	}

	balance := big.NewInt(0)
	for _, t := range mv.Transfers {
		balance.Add(balance, t.Value())
	}

	if expectedBalance.Cmp(balance) != 0 {
		t.Fatalf("expected balance is %s but got %s", expectedBalance.String(), balance.String())
	}
}

func TestGetAccountMovements_Ethereum(t *testing.T) {
	address := "0x7EF5A6135f1FD6a02593eEdC869c6D41D934aef8"
	since := uint64(8676237)
	expectedTransfers := []*domain.Transfer{
		{
			Type:        domain.Received,
			BlockHeight: 8676237,
			Amount:      big.NewInt(369136800000001),
		},
		{
			Type:        domain.Spent,
			BlockHeight: 8676237,
			Amount:      big.NewInt(1476547215001),
		},
		{
			Type:        domain.Received,
			BlockHeight: 8676237,
			Amount:      big.NewInt(0),
		},
		{
			Type:        domain.Spent,
			BlockHeight: 8676239,
			Amount:      big.NewInt(3152535117001),
		},
	}
	expectedBalance := big.NewInt(0)
	for _, t := range expectedTransfers {
		expectedBalance.Add(expectedBalance, t.Value())
	}

	api := blockbook.NewAPI(ethereumHostURL, blockbook.EthereumTranslator{})
	mv, err := api.GetAccountMovements(address, since)
	if err != nil {
		t.Fatal(err)
	}

	if len(mv.Transfers) != len(expectedTransfers) {
		t.Fatalf("expected to have %d changes but got %d changes", len(expectedTransfers), len(mv.Transfers))
	}

	balance := big.NewInt(0)
	for _, t := range mv.Transfers {
		balance.Add(balance, t.Value())
	}

	if expectedBalance.Cmp(balance) != 0 {
		t.Fatalf("expected balance is %s but got %s", expectedBalance.String(), balance.String())
	}
}

func TestGetLatestBlockHeight(t *testing.T) {
	api := blockbook.NewAPI(bitcoinHostURL, blockbook.BitcoinTranslator{})
	bh, err := api.GetLatestBlockHeight()
	if err != nil {
		t.Fatal(err)
	}

	if bh == 0 {
		t.Fatal("expected anything but 0")
	}
}
