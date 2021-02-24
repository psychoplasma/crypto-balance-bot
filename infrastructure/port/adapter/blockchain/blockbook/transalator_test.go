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

	if len(mvs.Transfers) != 1 {
		t.Fatalf("expected movement's balance change count is %d but got %d", 1, len(mvs.Transfers))
	}

	balanceDiff := big.NewInt(0)
	for _, t := range mvs.Transfers {
		balanceDiff = balanceDiff.Add(balanceDiff, t.Value())
	}

	if new(big.Int).Abs(balanceDiff).Uint64() == -value {
		t.Fatalf("expected movement's total balance change is %d but got %s", -value, balanceDiff.String())
	}
}
