package notification

import (
	"fmt"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

type MovementFormatter struct{}

func (mf MovementFormatter) Format(i interface{}) string {
	movementMap, _ := i.(map[*domain.Account][]*domain.AccountMovement)

	if len(movementMap) == 0 {
		return ""
	}

	msg := "```\n"
	for a, ms := range movementMap {
		mvmsg := ""
		for _, m := range ms {
			chmsg := ""
			for _, c := range m.Changes {
				chmsg += fmt.Sprintf(" => %s", c.Amount.String())
			}
			mvmsg += fmt.Sprintf("\tblock#%d{%s}\n", m.BlockHeight, chmsg)
		}
		msg += fmt.Sprintf("%s[%s]\n{\n%s}\n", a.Currency().Symbol, a.Address(), mvmsg)
	}
	msg += "```"

	return msg
}
