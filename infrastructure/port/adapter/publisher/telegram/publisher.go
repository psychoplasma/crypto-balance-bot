package telegram

import (
	"log"

	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/publisher"
	telegram "gopkg.in/tucnak/telebot.v2"
)

type recipient string

func (r recipient) Recipient() string {
	return string(r)
}

// Publisher is a message publisher using Telegram API
type Publisher struct {
	teleBot   *telegram.Bot
	formatter publisher.Formatter
}

// NewPublisher creates a new instance of TelegramPublisher
func NewPublisher(token string, fmt publisher.Formatter) *Publisher {
	bot, err := telegram.NewBot(telegram.Settings{
		Token: token,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Publisher{
		teleBot:   bot,
		formatter: fmt,
	}
}

// PublishMessage sends the given message to the given telegram user
func (tp *Publisher) PublishMessage(userID string, i interface{}) {
	msg := tp.formatter(i)

	if msg == "" {
		return
	}

	if _, err := tp.teleBot.Send(recipient(userID), msg, telegram.ModeMarkdown); err != nil {
		log.Printf("failed to send message: %s to telegram user: %s\n", msg, userID)
	}
}
