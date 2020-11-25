package blockchaindotcom_test

import (
	"math/big"
	"testing"

	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchaindotcom"
)

func TestToAccountMovements(t *testing.T) {
	addr1 := "test-addr-1"
	addr2 := "test-addr-2"
	txs := []blockchaindotcom.Transaction{
		{
			BlockHeight: 10,
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

	mvs := tr.ToAccountMovements(addr1, txs)

	if len(mvs) != 1 {
		t.Fatalf("expected movements count is %d but got %d", 1, len(mvs))
	}

	if mvs[0].BlockHeight != 10 {
		t.Fatalf("expected movement's block height is %d but got %d", 10, mvs[0].BlockHeight)
	}

	if len(mvs[0].Changes) != 2 {
		t.Fatalf("expected movement's balance change count is %d but got %d", 2, len(mvs[0].Changes))
	}

	balanceDiff := big.NewInt(0)
	for _, ch := range mvs[0].Changes {
		balanceDiff = balanceDiff.Add(balanceDiff, ch.Amount)
	}

	if balanceDiff.Cmp(big.NewInt(-2)) != 0 {
		t.Fatalf("expected movement's total balance change is %d but got %s", -2, balanceDiff.String())
	}
}
