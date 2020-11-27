package telegram

import (
	"fmt"
	"math/big"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

// MovementFormatter formats the given account movements to a string representation for telegram publisher
func MovementFormatter(v interface{}) string {
	acms, _ := v.([]*domain.AccountMovement)

	if len(acms) == 0 {
		return ""
	}

	movementExist := false
	for _, am := range acms {
		if len(am.Changes) > 0 {
			movementExist = true
		}
	}

	if !movementExist {
		return ""
	}

	msg := "```\n"
	for _, am := range acms {
		mvmsg := ""
		for blockHeight, chs := range am.Changes {
			chmsg := ""
			for _, c := range chs {
				// c.Amount / am.Currency.Decimal with 6 floating precision.
				// For example amount:5, symbol:eth, decimal: 1000
				// then the resulting string would be " => 0.00500 eth"
				chmsg += fmt.Sprintf(" => %s %s",
					new(big.Float).Quo(new(big.Float).SetInt(c.Amount),
						new(big.Float).SetInt(am.Currency.Decimal)).Text('f', 6),
					am.Currency.Symbol)
			}
			mvmsg += fmt.Sprintf("\tblock#%d{%s }\n", blockHeight, chmsg)
		}
		msg += fmt.Sprintf("%s[%s]\n{\n%s}\n", am.Currency.Symbol, am.Address, mvmsg)
	}
	msg += "```"

	return msg
}
