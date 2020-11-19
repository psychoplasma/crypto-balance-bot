package application

type notifier func(recipient string, msg interface{})

type Formatter interface {
	formatMovement(msg interface{}) (string, error)
	formatValue(msg interface{}) (string, error)
}

var formatter Formatter

func check() {

}

// func checkMovements(s *cryptoBot.Subscription, n notifier) {
// 	movements := s.Accounts.UpdateTxs()
// 	if len(movements) < 1 {
// 		return
// 	}

// 	msg, err := formatter.formatMovement(movements)
// 	if err != nil {
// 		log.Printf("cannot format account movements for user %s", s.UserID)
// 		return
// 	}

// 	n(s.UserID, msg)
// }

// func checkValue(s *cryptoBot.Subscription, n notifier) {
// 	movements := s.Account.UpdateBalances()
// 	if len(movements) < 1 {
// 		return
// 	}

// 	msg, err := formatter.formatValue(movements)
// 	if err != nil {
// 		log.Printf("cannot format account movements for user %s", s.UserID)
// 		return
// 	}

// 	n(s.UserID, msg)
// }
