// +integration

package blockchaindotcom_test

import (
	"testing"

	cryptobot "github.com/psychoplasma/crypto-balance-bot"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/notification"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/blockchaindotcom"
)

func TestGetTxsOfAddress(t *testing.T) {
	api := blockchaindotcom.NewBitcoinAPI(blockchaindotcom.BitcoinTranslator{})

	mv, err := api.GetTxsOfAddress("1AJbsFZ64EpEfS5UAjAfcUG8pH8Jn3rn1F", 157240)
	if err != nil {
		t.Fatal(err)
	}

	f := notification.MovementFormatter{}
	t.Log(f.Format(map[*cryptobot.Account][]*cryptobot.AccountMovement{
		cryptobot.NewAccount(cryptobot.Currency{}, "1AJbsFZ64EpEfS5UAjAfcUG8pH8Jn3rn1F"): mv,
	}))

	chagesExistForBlock := false
	for _, m := range mv {
		chagesExistForBlock = m.BlockHeight == 157240
	}

	if !chagesExistForBlock {
		t.Fatalf("expected to have changes in block#%d but got nothing", 157240)
	}
}
