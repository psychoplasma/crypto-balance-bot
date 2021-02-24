package etherscanio_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchain/etherscanio"
)

func TestToAccountMovements(t *testing.T) {
	addr1 := "test-addr-1"
	addr1MixedCase := "teSt-addR-1"
	addr2 := "test-addr-2"
	blockHeight := uint64(10)
	txs := []etherscanio.Transaction{
		{
			BlockHeight: fmt.Sprint(blockHeight),
			From:        addr1MixedCase,
			To:          addr2,
			Value:       "100",
			Status:      "1",
			Timestamp:   "1610503881",
		},
		{
			BlockHeight: fmt.Sprint(blockHeight),
			From:        addr1,
			To:          addr2,
			Value:       "100",
			Status:      "0",
			Timestamp:   "1610503881",
		},
	}

	mv, err := etherscanio.EthereumTranslator{}.ToAccountMovements(addr1, txs)
	if err != nil {
		t.Fatal(err)
	}

	if len(mv.Transfers) != 1 {
		t.Fatalf("expected movements count is %d but got %d", 1, len(mv.Transfers))
	}

	balanceDiff := new(big.Int)
	for _, t := range mv.Transfers {
		balanceDiff = balanceDiff.Add(balanceDiff, t.Value())
	}

	if balanceDiff.Cmp(big.NewInt(-100)) != 0 {
		t.Fatalf("expected movement's total balance change is %d but got %s", -100, balanceDiff.String())
	}
}
