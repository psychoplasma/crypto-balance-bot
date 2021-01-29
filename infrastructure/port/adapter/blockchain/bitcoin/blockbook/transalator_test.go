package blockbook_test

import (
	"math/big"
	"testing"

	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchain/bitcoin/blockchaindotcom"
)

func TestToAccountMovements(t *testing.T) {
	addr1 := "test-addr-1"
	addr2 := "test-addr-2"
	blockHeight := 10
	txs := []blockchaindotcom.Transaction{
		{
			BlockHeight: blockHeight,
			Hash:        "hash1",
			Inputs: []blockchaindotcom.Input{
				{
					PrevOutput: blockchaindotcom.Output{
						Address: addr1,
						Value:   big.NewInt(5),
					},
				},
				{
					PrevOutput: blockchaindotcom.Output{
						Address: addr2,
						Value:   big.NewInt(3),
					},
				},
			},
			Outputs: []blockchaindotcom.Output{
				{
					Address: addr1,
					Value:   big.NewInt(3),
				},
				{
					Address: addr2,
					Value:   big.NewInt(5),
				},
			},
		},
	}

	tr := blockchaindotcom.BitcoinTranslator{}

	mvs, err := tr.ToAccountMovements(addr1, txs)
	if err != nil {
		t.Fatal(err)
	}

	if len(mvs.Changes) != 1 {
		t.Fatalf("expected movements count is %d but got %d", 1, len(mvs.Changes))
	}

	if mvs.Changes[blockHeight] == nil {
		t.Fatalf("expected to have changes at block#%d but got nothing", blockHeight)
	}

	if len(mvs.Changes[blockHeight]) != 2 {
		t.Fatalf("expected movement's balance change count is %d but got %d", 2, len(mvs.Changes[blockHeight]))
	}

	balanceDiff := big.NewInt(0)
	for _, ch := range mvs.Changes[blockHeight] {
		balanceDiff = balanceDiff.Add(balanceDiff, ch.Amount)
	}

	if balanceDiff.Cmp(big.NewInt(-2)) != 0 {
		t.Fatalf("expected movement's total balance change is %d but got %s", -2, balanceDiff.String())
	}
}
