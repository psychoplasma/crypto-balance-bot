package blockchaindotcom_test

import (
	"math/big"
	"testing"

	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchain/blockchaindotcom"
)

func TestToAccountMovements(t *testing.T) {
	addr1 := "test-addr-1"
	addr2 := "test-addr-2"
	blockHeight := uint64(10)
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

	if len(mvs.Transfers) != 2 {
		t.Fatalf("expected movement's balance change count is %d but got %d", 2, len(mvs.Transfers))
	}

	balanceDiff := big.NewInt(0)
	for _, t := range mvs.Transfers {
		balanceDiff = balanceDiff.Add(balanceDiff, t.Value())
	}

	if balanceDiff.Cmp(big.NewInt(-2)) != 0 {
		t.Fatalf("expected movement's total balance change is %d but got %s", -2, balanceDiff.String())
	}
}
