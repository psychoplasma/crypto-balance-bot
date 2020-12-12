package telegram

import (
	"fmt"
	"math/big"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

// MovementFormatter formats the given account movements to a string representation for telegram publisher
func MovementFormatter(v interface{}) string {
	event, _ := v.(*domain.AccountAssetsMovedEvent)
	acms := event.AccountMovements()
	currency := event.Currency()

	// We don't want to create empty movement message
	// Instead setting the telegram message to empty string
	// will make telegram bot not to the send any message at all
	if !doesMovementExist(acms) {
		return ""
	}

	// Format the message as follows:
	// ```
	// symbol[address]
	// {
	//   block#n{ => amount symbol => ... }
	//   block#n+1{ => amount symbol => ... }
	//   .
	//   .
	//   block#m{ => amount symbol => ... }
	// }
	// ```
	msg := ""
	for _, blockHeight := range acms.Blocks {
		chmsg := ""
		for _, c := range acms.Changes[blockHeight] {
			// c.Amount / am.Currency.Decimal with 6 floating precision.
			// For example amount:5, symbol:eth, decimal: 1000
			// then the resulting string would be " => 0.00500 eth"
			chmsg += fmt.Sprintf(" => %s %s",
				new(big.Float).Quo(new(big.Float).SetInt(c.Amount),
					new(big.Float).SetInt(currency.Decimal)).Text('f', 6),
				currency.Symbol)
		}
		msg += fmt.Sprintf("\tblock#%d{%s }\n", blockHeight, chmsg)
	}
	msg = fmt.Sprintf("```\n%s[%s]\n{\n%s}\n```", currency.Symbol, acms.Address, msg)

	return msg
}

func doesMovementExist(acms *domain.AccountMovements) bool {
	for _, ch := range acms.Changes {
		if len(ch) > 0 {
			return true
		}
	}

	return false
}
