package notification

import (
	"fmt"
	"log"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

type MovementFormatter struct{}

func (mf MovementFormatter) Format(i interface{}) string {
	movementMap, s := i.(map[*domain.Account][]*domain.AccountMovement)
	if !s {
		log.Printf("cannot convert message: %#v to map[string][]*domain.AccountMovement\n", i)
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
