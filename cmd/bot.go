package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/psychoplasma/crypto-balance-bot/cmd/application"
	"github.com/psychoplasma/crypto-balance-bot/repo/inmemory"
	tb "gopkg.in/tucnak/telebot.v2"
	"gopkg.in/yaml.v2"
)

const parameterSeparator = " "

type botCommand struct {
	Usage          string
	Description    string
	ParameterCount int
}

var commands = map[string]botCommand{
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

var subsRepo = inmemory.NewSubscriptionReposititory()
var subsAppService = application.NewSubscriptionService(subsRepo)
var currencyAppService = application.NewCurrencyService()

func main() {
	startBot()
}

func startBot() {
	c, err := readConfig("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  c.Token,
		Poller: &tb.LongPoller{Timeout: c.PollingTime * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	b.Handle("/subscribe-for-value", func(m *tb.Message) {
		fmt.Printf("command: %#v\n", m)
		msg := strings.Split(m.Payload, parameterSeparator)
		fmt.Printf("message payload: %#v\n", msg)

		if len(msg) != commands["subscribe_for_value"].ParameterCount {
			b.Send(m.Sender, fmt.Sprintf("Invalid inputs, see command usage: \"%s\"", commands["subscribe_for_value"].Usage))
			return
		}

		if err := subsAppService.SubscribeForValue(
			m.Sender.Recipient(),
			msg[0],
			*currencyAppService.GetCurrency(msg[1]),
			msg[2],
			*currencyAppService.GetCurrency(msg[3]),
		); err != nil {
			log.Printf("failed to subscribe for value, %s", err.Error())
		}

		b.Send(m.Sender, fmt.Sprintf("subscribed %s:%s for", msg[0], msg[1]))
	})

	b.Handle("/subscribe_for_movement", func(m *tb.Message) {
		fmt.Printf("command: %#v\n", m)
		msg := strings.Split(m.Payload, parameterSeparator)
		fmt.Printf("message payload: %#v\n", msg)

		if len(msg) != commands["subscribe_for_movement"].ParameterCount {
			b.Send(m.Sender, fmt.Sprintf("Invalid inputs, see command usage: \"%s\"", commands["subscribe_for_movement"].Usage))
			return
		}

		if err := subsAppService.SubscribeForMovement(
			m.Sender.Recipient(),
			msg[0],
			*currencyAppService.GetCurrency(msg[1]),
			msg[2],
		); err != nil {
			log.Printf("failed to subscribe for value, %s", err.Error())
		}

		b.Send(m.Sender, fmt.Sprintf("subscribed %s:%s for", msg[0], msg[1]))
	})

	b.Handle("/unsubscribe", func(m *tb.Message) {
		fmt.Printf("message payload: %#v\n", m)

		if err := subsAppService.Unsubscribe(m.Payload); err != nil {
			b.Send(m.Sender, fmt.Sprintf("failed to unsubscribe, %s", err.Error()))
		}
	})

	b.Handle("/my_subscriptions", func(m *tb.Message) {
		subs, err := subsAppService.GetSubscriptionsForUser(m.Sender.Recipient())
		if err != nil {
			log.Printf("failed to subscribe for value, %s", err.Error())
			return
		}

		subsMsg := ""

		for _, s := range subs {
			addrs := ""
			for _, a := range s.Account.AddressList {
				addrs += " : " + a
			}
			subsMsg += fmt.Sprintf("ID: `%s`, Currency: `%s`, Addresses: `%s` \n", s.ID, s.Account.Currency.Symbol, addrs)
		}

		b.Send(m.Sender, subsMsg)
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		fmt.Printf("Unhandled message: %#v\n", m.Payload)
	})

	b.Start()
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
