package telegram

import (
	"fmt"
	"math/big"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

// MovementFormatter formats the given account movements to a string representation for telegram publisher
func MovementFormatter(v interface{}) string {
	sm, _ := v.(*domain.SubscriptionMovements)
	acms := sm.AccountMovements()

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
	msg := "```\n"
	for _, am := range acms {
		mvmsg := ""
		for _, blockHeight := range am.Blocks {
			chmsg := ""
			for _, c := range am.Changes[blockHeight] {
				// c.Amount / am.Currency.Decimal with 6 floating precision.
				// For example amount:5, symbol:eth, decimal: 1000
				// then the resulting string would be " => 0.00500 eth"
				chmsg += fmt.Sprintf(" => %s %s",
					new(big.Float).Quo(new(big.Float).SetInt(c.Amount),
						new(big.Float).SetInt(sm.Currency().Decimal)).Text('f', 6),
					sm.Currency().Symbol)
			}
			mvmsg += fmt.Sprintf("\tblock#%d{%s }\n", blockHeight, chmsg)
		}
		msg += fmt.Sprintf("%s[%s]\n{\n%s}\n", sm.Currency().Symbol, am.Address, mvmsg)
	}
	msg += "```"

	return msg
}

func doesMovementExist(acms map[string]*domain.AccountMovements) bool {
	movementExist := false
	for _, acm := range acms {
		for _, ch := range acm.Changes {
			if len(ch) > 0 {
				movementExist = true
				break
			}
		}

		if movementExist {
			break
		}
	}

	return movementExist
}
