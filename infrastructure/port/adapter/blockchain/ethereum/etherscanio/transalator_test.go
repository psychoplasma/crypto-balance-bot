package etherscanio_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchain/ethereum/etherscanio"
)

func TestToAccountMovements(t *testing.T) {
	addr1 := "test-addr-1"
	addr1MixedCase := "teSt-aDdR-1"
	addr2 := "test-addr-2"
	blockHeight := 10
	txs := []etherscanio.Transaction{
		{
			BlockHeight: fmt.Sprint(blockHeight),
			From:        addr1MixedCase,
			To:          addr2,
			Value:       "100",
			Status:      "1",
		},
		{
			BlockHeight: fmt.Sprint(blockHeight),
			From:        addr1,
			To:          addr2,
			Value:       "100",
			Status:      "0",
		},
	}

	mv, err := etherscanio.EthereumTranslator{}.ToAccountMovements(addr1, txs)
	if err != nil {
		t.Fatal(err)
	}

	if len(mv.Changes) != 1 {
		t.Fatalf("expected movements count is %d but got %d", 1, len(mv.Changes))
	}

	if mv.Changes[blockHeight] == nil {
		t.Fatalf("expected to have balance changes at block#%d but got nothing", blockHeight)
	}

	if len(mv.Changes[blockHeight]) != 1 {
		t.Fatalf("expected balance change count is %d but got %d", 2, len(mv.Changes[blockHeight]))
	}

	balanceDiff := new(big.Int)
	for _, ch := range mv.Changes[blockHeight] {
		balanceDiff = balanceDiff.Add(balanceDiff, ch.Amount)
	}

	if balanceDiff.Cmp(big.NewInt(-100)) != 0 {
		t.Fatalf("expected movement's total balance change is %d but got %s", -100, balanceDiff.String())
	}
}
