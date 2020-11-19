package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/psychoplasma/crypto-balance-bot/application"
	tb "gopkg.in/tucnak/telebot.v2"
	"gopkg.in/yaml.v2"
)

const parameterSeparator = " "

var commands = map[string]command{
	"subscribe_for_movement": {
		Usage:          "/subscribe_for_movement <name> <ticker> <address descriptor>",
		Description:    "",
		ParameterCount: 3,
	},
	"subscribe_for_value": {
		Usage:          "/subscribe_for_value <name> <ticker> <address descriptor> <against ticker>",
		Description:    "",
		ParameterCount: 4,
	},
	"unsubscribe": {
		Usage:          "/unsubscribe <subscription id>",
		Description:    "",
		ParameterCount: 1,
	},
	"my_subscriptions": {
		Usage:          "/my_subscriptions",
		Description:    "",
		ParameterCount: 0,
	},
}

type config struct {
	Token       string        `yaml:"token"`
	PollingTime time.Duration `yaml:"polling-time"`
}

type command struct {
	Usage          string
	Description    string
	ParameterCount int
}

type Bot struct {
	tb          *tb.Bot
	subsApp     *application.SubscriptionApplication
	currencyApp *application.CurrencyService
}

func NewBot(subsApp *application.SubscriptionApplication, currencyApp *application.CurrencyService) Bot {
	c, err := readConfig("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tb.NewBot(tb.Settings{
		Token:  c.Token,
		Poller: &tb.LongPoller{Timeout: c.PollingTime * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	return Bot{
		tb:          bot,
		subsApp:     subsApp,
		currencyApp: currencyApp,
	}
}

func (b Bot) RegisterCommands() {
	b.tb.Handle("/subscribe_for_value", func(m *tb.Message) {
		fmt.Printf("command: %#v\n", m)
		msg := strings.Split(m.Payload, parameterSeparator)
		fmt.Printf("message payload: %#v\n", msg)

		if len(msg) != commands["subscribe_for_value"].ParameterCount {
			b.tb.Send(m.Sender, fmt.Sprintf("Invalid inputs, see command usage: \"%s\"", commands["subscribe_for_value"].Usage))
			return
		}

		if err := b.subsApp.SubscribeForValue(
			m.Sender.Recipient(),
			msg[0],
			*b.currencyApp.GetCurrency(msg[1]),
			msg[2],
			*b.currencyApp.GetCurrency(msg[3]),
		); err != nil {
			log.Printf("failed to subscribe for value, %s", err.Error())
		}

		b.tb.Send(m.Sender, fmt.Sprintf("subscribed %s:%s for", msg[0], msg[1]))
	})

	b.tb.Handle("/subscribe_for_movement", func(m *tb.Message) {
		fmt.Printf("command: %#v\n", m)
		msg := strings.Split(m.Payload, parameterSeparator)
		fmt.Printf("message payload: %#v\n", msg)

		if len(msg) != commands["subscribe_for_movement"].ParameterCount {
			b.tb.Send(m.Sender, fmt.Sprintf("Invalid inputs, see command usage: \"%s\"", commands["subscribe_for_movement"].Usage))
			return
		}

		if err := b.subsApp.SubscribeForMovement(
			m.Sender.Recipient(),
			msg[0],
			*b.currencyApp.GetCurrency(msg[1]),
			msg[2],
		); err != nil {
			log.Printf("failed to subscribe for value, %s", err.Error())
		}

		b.tb.Send(m.Sender, fmt.Sprintf("subscribed %s:%s for", msg[0], msg[1]))
	})

	b.tb.Handle("/unsubscribe", func(m *tb.Message) {
		fmt.Printf("message payload: %#v\n", m)

		if err := b.subsApp.Unsubscribe(m.Payload); err != nil {
			b.tb.Send(m.Sender, fmt.Sprintf("failed to unsubscribe, %s", err.Error()))
		}
	})

	b.tb.Handle("/my_subscriptions", func(m *tb.Message) {
		subs, err := b.subsApp.GetSubscriptionsForUser(m.Sender.Recipient())
		if err != nil {
			log.Printf("failed to subscribe for value, %s", err.Error())
			return
		}

		subsMsg := ""

		log.Printf("Subscription count: %d\n", len(subs))
		for _, s := range subs {
			addrs := ""
			log.Printf("Address count: %d\n", len(addrs))
			for _, a := range s.Accounts {
				addrs += " : " + a.Address()
			}
			subsMsg += fmt.Sprintf("ID: `%s`, Addresses: `%s` \n", s.ID, addrs)
		}

		log.Println(subsMsg)

		b.tb.Send(m.Sender, subsMsg)
	})

	b.tb.Handle(tb.OnText, func(m *tb.Message) {
		fmt.Printf("Unhandled message: %#v\n", m.Payload)
	})
}

func (b Bot) Start() {
	b.tb.Start()
}

func readConfig(path string) (*config, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, fmt.Errorf("configuration file does not exist")
	}

	c := &config{}
	err = yaml.Unmarshal(f, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
