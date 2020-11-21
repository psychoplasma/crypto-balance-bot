package adapter

import (
	"log"

	"github.com/psychoplasma/crypto-balance-bot/infrastructure/notification"
	telegram "gopkg.in/tucnak/telebot.v2"
)

type telegramRecipient string

func (tr telegramRecipient) Recipient() string {
	return string(tr)
}

// TelegramPublisher is a message publisher using Telegram API
type TelegramPublisher struct {
	teleBot *telegram.Bot
	fmt     notification.Formatter
}

// NewTelegramPublisher creates a new instance of TelegramPublisher
func NewTelegramPublisher(token string, fmt notification.Formatter) *TelegramPublisher {
	bot, err := telegram.NewBot(telegram.Settings{
		Token: token,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &TelegramPublisher{
		teleBot: bot,
		fmt:     fmt,
	}
}

// PublishMessage sends the given message to the given telegram user
func (tp *TelegramPublisher) PublishMessage(userID string, i interface{}) {
	msg := tp.fmt.Format(i)
	if _, err := tp.teleBot.Send(telegramRecipient(userID), msg, telegram.ModeMarkdown); err != nil {
		log.Printf("failed to send message: %s to telegram user: %s\n", msg, userID)
	}
}
