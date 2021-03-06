package blockbook_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchain/blockbook"
)

func TestBitcoinTranslator_ToAccountMovements(t *testing.T) {
	addr1 := "test-addr-1"
	addr2 := "test-addr-2"
	blockHeight := uint64(10)
	addrTxs := []blockbook.Transaction{
		{
			BlockHeight: blockHeight,
			Inputs: []blockbook.Input{
				{
					Addresses: []string{addr1},
					Value:     "5",
				},
				{
					Addresses: []string{addr2},
					Value:     "3",
				},
			},
			Outputs: []blockbook.Output{
				{
					Addresses: []string{addr1},
					Value:     "3",
				},
				{
					Addresses: []string{addr2},
					Value:     "5",
				},
			},
		},
	}

	mvs, err := new(blockbook.BitcoinTranslator).ToAccountMovements(addr1, addrTxs)
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
		balanceDiff = balanceDiff.Add(balanceDiff, ch.Value())
	}

	if balanceDiff.Cmp(big.NewInt(-2)) != 0 {
		t.Fatalf("expected movement's total balance change is %d but got %s", -2, balanceDiff.String())
	}
}

func TestEthereumTranslator_ToAccountMovements(t *testing.T) {
	addr1 := "test-addr-1"
	addr2 := "test-addr-2"
	blockHeight := uint64(10)
	value := uint64(5)
	addrTxs := []blockbook.Transaction{
		{
			BlockHeight: blockHeight,
			Value:       fmt.Sprint(value),
			Inputs: []blockbook.Input{
				{
					Addresses: []string{addr1},
				},
			},
			Outputs: []blockbook.Output{
				{
					Addresses: []string{addr2},
				},
			},
			EthereumSpecific: blockbook.EthereumSpecific{
				Status: 1,
			},
		},
	}

	mvs, err := new(blockbook.EthereumTranslator).ToAccountMovements(addr1, addrTxs)
	if err != nil {
		t.Fatal(err)
	}

	if len(mvs.Changes) != 1 {
		t.Fatalf("expected movements count is %d but got %d", 1, len(mvs.Changes))
	}

	if mvs.Changes[blockHeight] == nil {
		t.Fatalf("expected to have changes at block#%d but got nothing", blockHeight)
	}

	if len(mvs.Changes[blockHeight]) != 1 {
		t.Fatalf("expected movement's balance change count is %d but got %d", 1, len(mvs.Changes[blockHeight]))
	}

	balanceDiff := big.NewInt(0)
	for _, ch := range mvs.Changes[blockHeight] {
		balanceDiff = balanceDiff.Add(balanceDiff, ch.Value())
	}

	if new(big.Int).Abs(balanceDiff).Uint64() == -value {
		t.Fatalf("expected movement's total balance change is %d but got %s", -value, balanceDiff.String())
	}
}
