package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/psychoplasma/crypto-balance-bot/application"
	tb "gopkg.in/tucnak/telebot.v2"
)

const parameterSeparator = " "

var commands = map[string]command{
	"subscribe_for_movement": {
		Usage:          "/subscribe_for_movement <name> <ticker> [<address descriptor> ...]",
		Description:    "",
		ParameterCount: 3,
	},
	"subscribe_for_value": {
		Usage:          "/subscribe_for_value <name> <ticker> <against ticker> <address descriptor> ",
		Description:    "",
		ParameterCount: 4,
	},
	"unsubscribe": {
		Usage:          "/unsubscribe <subscription id>",
		Description:    "",
		ParameterCount: 1,
	},
	"unsubscribe_all": {
		Usage:          "/unsubscribe_all",
		Description:    "",
		ParameterCount: 0,
	},
	"my_subscriptions": {
		Usage:          "/my_subscriptions",
		Description:    "",
		ParameterCount: 0,
	},
}

type command struct {
	Usage          string
	Description    string
	ParameterCount int
}

// Bot is TelegramBot receives subscription related commands from a user and returns the corresponding responses
type Bot struct {
	tb      *tb.Bot
	subsApp *application.SubscriptionApplication
}

// NewBot creates a new instance of Bot
func NewBot(c *Config, subsApp *application.SubscriptionApplication) Bot {
	bot, err := tb.NewBot(tb.Settings{
		Token:  c.Token,
		Poller: &tb.LongPoller{Timeout: c.PollingTime * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	return Bot{
		tb:      bot,
		subsApp: subsApp,
	}
}

// Start starts the bot
func (b Bot) Start() {
	b.registerCommands()
	b.tb.Start()
}

// Stop stops the bot gracefully
func (b Bot) Stop() {
	b.tb.Stop()
}

func (b Bot) registerCommands() {
	b.tb.Handle("/subscribe_for_value", func(m *tb.Message) {
		msg := strings.Split(m.Payload, parameterSeparator)
		fmt.Printf("message payload: %#v\n", msg)

		if len(msg) < commands["subscribe_for_value"].ParameterCount {
			b.tb.Send(m.Sender,
				fmt.Sprintf("Invalid inputs, see command usage: \"%s\"",
					commands["subscribe_for_value"].Usage))
			return
		}

		if err := b.subsApp.SubscribeForValue(
			m.Sender.Recipient(),
			msg[0],
			msg[1],
			msg[2],
			msg[3:],
		); err != nil {
			log.Printf("failed to subscribe for value, %s", err.Error())
		}

		b.tb.Send(m.Sender, fmt.Sprintf("subscribed %s:%s for", msg[0], msg[1]))
	})

	b.tb.Handle("/subscribe_for_movement", func(m *tb.Message) {
		msg := strings.Split(m.Payload, parameterSeparator)
		fmt.Printf("message payload: %#v\n", msg)

		if len(msg) < commands["subscribe_for_movement"].ParameterCount {
			b.tb.Send(m.Sender, fmt.Sprintf("Invalid inputs, see command usage: \"%s\"", commands["subscribe_for_movement"].Usage))
			return
		}

		if err := b.subsApp.SubscribeForMovement(
			m.Sender.Recipient(),
			msg[0],
			msg[1],
			msg[2:],
		); err != nil {
			log.Printf("failed to subscribe for value, %s", err.Error())
		}

		b.tb.Send(m.Sender, fmt.Sprintf("subscribed %s:%s for", msg[0], msg[1]))
	})

	b.tb.Handle("/unsubscribe", func(m *tb.Message) {
		if err := b.subsApp.Unsubscribe(m.Payload); err != nil {
			b.tb.Send(m.Sender, fmt.Sprintf("failed to unsubscribe, %s", err.Error()))
		}
	})

	b.tb.Handle("/unsubscribe_all", func(m *tb.Message) {
		if err := b.subsApp.UnsubscribeAllForUser(m.Sender.Recipient()); err != nil {
			b.tb.Send(m.Sender, fmt.Sprintf("failed to unsubscribe all, %s", err.Error()))
		}
	})

	b.tb.Handle("/my_subscriptions", func(m *tb.Message) {
		subs, err := b.subsApp.GetSubscriptionsForUser(m.Sender.Recipient())
		if err != nil {
			log.Printf("failed to subscribe for value, %s", err.Error())
			return
		}

		subsMsg := ""

		for _, s := range subs {
			addrs := ""
			log.Printf("Address count: %d\n", len(addrs))
			for _, a := range s.Accounts() {
				addrs += " : " + a.Address()
			}
			subsMsg += fmt.Sprintf("ID: %s, Addresses: %s \n", s.ID(), addrs)
		}

		log.Println(subsMsg)

		b.tb.Send(m.Sender, subsMsg)
	})

	b.tb.Handle(tb.OnText, func(m *tb.Message) {
		fmt.Printf("Unhandled message: %#v\n", m.Payload)
	})
}
