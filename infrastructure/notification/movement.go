package notification

import (
	"fmt"
	"math/big"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

// MovementFormatter formats the given account movements to a string representation for telegram publisher
func MovementFormatter(i interface{}) string {
	movementMap, _ := i.(map[*domain.Account][]*domain.AccountMovement)

	if len(movementMap) == 0 {
		return ""
	}

	// FIXME: this is also inefficient
	totalMovements := 0
	for _, ms := range movementMap {
		totalMovements += len(ms)
	}

	if totalMovements == 0 {
		return ""
	}

	msg := "```\n"
	for a, ms := range movementMap {
		mvmsg := ""
		for _, m := range ms {
			chmsg := ""
			for _, c := range m.Changes {
				// FIXME: this looks too inefficient
				chmsg += fmt.Sprintf(" => %s %s",
					big.NewFloat(0).Quo(new(big.Float).SetInt(c.Amount), new(big.Float).SetInt(a.Currency().Decimal)).Text('f', 6),
					a.Currency().Symbol)
			}
			mvmsg += fmt.Sprintf("\tblock#%d{%s}\n", m.BlockHeight, chmsg)
		}
		msg += fmt.Sprintf("%s[%s]\n{\n%s}\n", a.Currency().Symbol, a.Address(), mvmsg)
	}
	msg += "```"

	return msg
}
