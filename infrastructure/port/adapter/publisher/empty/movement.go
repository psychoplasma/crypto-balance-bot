package empty

import (
	"fmt"
	"math/big"
	"time"

	domain "github.com/psychoplasma/crypto-balance-bot"
)

// MovementFormatter formats the given account movements to a string representation for telegram publisher
func MovementFormatter(v interface{}) string {
	event, _ := v.(*domain.AccountAssetsMovedEvent)
	transfers := event.Transfers()
	account := event.Account()
	currency := event.Currency()

	// We don't want to create empty movement message
	// Instead setting the telegram message to empty string
	// will make telegram bot not to the send any message at all
	if len(transfers) == 0 {
		return ""
	}

	// Format the message as follows:
	// ```
	// <address> Received
	// {
	//   <from address>
	//   <amount> <symbol>
	//   <time>
	//   <block#>
	// }
	// <address> Spent
	// {
	//   <to address>
	//   <amount> <symbol>
	//   <time>
	//   <block#>
	// }
	// ```
	// <from address> and <to address> are applicable only for account-based blockchains
	msg := ""
	for _, t := range transfers {
		switch t.Type {
		case domain.Received:
			msg += fmt.Sprintf("%s Received\n{\n", account)
			break
		case domain.Spent:
			msg += fmt.Sprintf("%s Spent\n{\n", account)
			break
		}

		msg += fmt.Sprintf("\t%s\n", t.Address)
		// c.Amount / am.Currency.Decimal with 6 floating precision.
		// For example amount:5, symbol:eth, decimal: 1000
		// then the resulting string would be " => 0.00500 eth"
		msg += fmt.Sprintf("\t%s %s\n",
			new(big.Float).Quo(new(big.Float).SetInt(t.Amount),
				new(big.Float).SetInt(currency.Decimal)).Text('f', 6),
			currency.Symbol)
		msg += fmt.Sprintf("\ttime@%s\n", time.Unix(int64(t.Timestamp), 0).Format(time.RFC3339))
		msg += fmt.Sprintf("\tblock#%d\n}\n", t.BlockHeight)
	}
	msg = fmt.Sprintf("```\n%s```", msg)

	return msg
}
