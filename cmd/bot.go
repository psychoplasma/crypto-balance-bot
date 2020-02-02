package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
	"gopkg.in/yaml.v2"
)

const parameterSeparator = " "

var subscriptions map[string]*Subscription

type config struct {
	Token       string        `yaml:"token"`
	PollingTime time.Duration `yaml:"polling-time"`
}

func main() {
	subscriptions := make(map[string]*Subscription, 1)

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
