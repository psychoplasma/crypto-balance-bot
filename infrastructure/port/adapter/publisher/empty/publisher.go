package empty

import (
	"log"

	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/publisher"
)

type recipient string

func (r recipient) Recipient() string {
	return string(r)
}

// Publisher is a message publisher using Telegram API
type Publisher struct {
	formatter publisher.Formatter
}

// NewPublisher creates a new instance of EmptyPublisher
func NewPublisher(token string, fmt publisher.Formatter) *Publisher {
	return &Publisher{
		formatter: fmt,
	}
}

// PublishMessage sends the given message to the given telegram user
func (tp *Publisher) PublishMessage(userID string, i interface{}) {
	msg := tp.formatter(i)

	if msg == "" {
		return
	}

	log.Printf(msg)
}
