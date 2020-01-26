package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	cryptobot "github.com/psychoplasma/crypto-balance-bot"
	tb "gopkg.in/tucnak/telebot.v2"
	"gopkg.in/yaml.v2"
)

const parameterSeparator = " "

var subscriptions map[string]*cryptobot.Subscription

func main() {
	subscriptions := make(map[string]*cryptobot.Subscription, 1)

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

	b.Handle("/subscribe", func(m *tb.Message) {
		fmt.Printf("%#v\n", m.Payload)
		msg := strings.Split(m.Payload, parameterSeparator)
		fmt.Printf("%#v\n", msg)
		b.Send(m.Sender, fmt.Sprintf("subscribed %s:%s for", msg[0], msg[1]))
	})

	b.Handle("/unsubscribe", func(m *tb.Message) {
		fmt.Printf("%#v\n", m)
		b.Send(m.Sender, fmt.Sprintf("Removed subscription: %s", m.Payload))
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		fmt.Printf("Unhandled message: %#v\n", m.Payload)
	})

	b.Start()
}

type config struct {
	Token       string        `yaml:"token"`
	PollingTime time.Duration `yaml:"polling-time"`
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

// Subscribe subscribes for value change or account movement events for the given account
func subscribe(stype cryptobot.SubscriptionType) error {
	switch stype {
	case cryptobot.Value:
		return subscribeForValue()
	case cryptobot.Movement:
		return subscribeForMovement()
	default:
		return cryptobot.ErrSubscriptionType
	}
}

// Unsubscribe unsubscribes for value change or account movement events for the given account
func unsubscribe(stype cryptobot.SubscriptionType) error {
	switch stype {
	case cryptobot.Value:
		return unsubscribeForValue()
	case cryptobot.Movement:
		return unsubscribeForMovement()
	default:
		return cryptobot.ErrSubscriptionType
	}
}

func subscribeForValue(currency string, address string, currencyAgainst string) error {
	_, ok := subscriptions[fmt.Sprintf("%s:%s:%s", currency, address, currencyAgainst)]
	if ok {
		return errors.New("already subscribed")
	}

	subscriptions[fmt.Sprintf("%s:%s:%s", currency, address, currencyAgainst)] = &cryptobot.Subscription

	return nil
}

func subscribeForMovement() error {
	return nil
}

func unsubscribeForValue() error {
	return nil
}

func unsubscribeForMovement() error {
	return nil
}
