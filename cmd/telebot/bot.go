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
	"subscription_details": {
		Endpoint:       "/subscription",
		Usage:          "/subscription <subscription ID>",
		Description:    "Displays the details of the given subscription",
		ParameterCount: 1,
	},
	"subscribe": {
		Endpoint:       "/subscribe",
		Usage:          "/subscribe <asset symbol> [<account's address> ...]",
		Description:    "Subscribes to asset movements for an account",
		ParameterCount: 2,
	},
	"unsubscribe": {
		Endpoint:       "/unsubscribe",
		Usage:          "/unsubscribe <subscription ID>",
		Description:    "Deletes the given subscription of the sender",
		ParameterCount: 1,
	},
	"unsubscribe_all": {
		Endpoint:       "/unsubscribe_all",
		Usage:          "/unsubscribe_all",
		Description:    "Deletes all subscriptions of the sender",
		ParameterCount: 0,
	},
	"my_subscriptions": {
		Endpoint:       "/my_subscriptions",
		Usage:          "/my_subscriptions",
		Description:    "Shows all subscriptions of the sender",
		ParameterCount: 0,
	},
	"available_assets": {
		Endpoint:       "/assets",
		Usage:          "/assets",
		Description:    "Shows available asset for a subscription",
		ParameterCount: 0,
	},
	"available_commands": {
		Endpoint:       "/commands",
		Usage:          "/commands",
		Description:    "Shows available commands",
		ParameterCount: 0,
	},
	"help": {
		Endpoint:       "/help",
		Usage:          "/help",
		Description:    "Shows this message",
		ParameterCount: 0,
	},
}

type command struct {
	Endpoint       string
	Usage          string
	Description    string
	ParameterCount int
}

func (cmd command) ValidateParameters(payload string) ([]string, error) {
	params := strings.Split(payload, parameterSeparator)
	log.Printf("Command parameters: %#v\n", params)

	if len(params) < cmd.ParameterCount {
		return nil, fmt.Errorf("Wrong number of inputs, command usage:\n\n```\n\t%s```", cmd.Usage)
	}

	return params, nil
}

// Bot is TelegramBot receives subscription related commands from a user and returns the corresponding responses
type Bot struct {
	tb      *tb.Bot
	subsApp *application.SubscriptionApplication
}

// NewBot creates a new instance of Bot
func NewBot(c *Config, subsApp *application.SubscriptionApplication) Bot {
	bot, err := tb.NewBot(tb.Settings{
		Token:  c.Telebot.Token,
		Poller: &tb.LongPoller{Timeout: c.Telebot.PollingTime * time.Second},
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
	log.Println("Starting Telegram Bot")
	b.registerCommands()
	b.tb.Start()
}

// Stop stops the bot gracefully
func (b Bot) Stop() {
	b.tb.Stop()
}

func (b Bot) registerCommands() {
	b.tb.Handle(commands["subscription_details"].Endpoint, b.subscriptionDetailsCMD)
	b.tb.Handle(commands["subscribe"].Endpoint, b.subscribeForMovementCMD)
	b.tb.Handle(commands["unsubscribe"].Endpoint, b.unsubscribeCMD)
	b.tb.Handle(commands["unsubscribe_all"].Endpoint, b.unsubscribeAllCMD)
	b.tb.Handle(commands["my_subscriptions"].Endpoint, b.mySubscriptionsCMD)
	b.tb.Handle(commands["available_assets"].Endpoint, b.availableAssetsCMD)
	b.tb.Handle(commands["available_commands"].Endpoint, b.availableCommandsCMD)
	b.tb.Handle(commands["help"].Endpoint, b.helpCMD)

	// Default handler for unhandled commands
	b.tb.Handle(tb.OnText, b.defaultHandler)
}

func (b Bot) subscriptionDetailsCMD(m *tb.Message) {
	params, err := commands["subscription_details"].ValidateParameters(m.Payload)
	if err != nil {
		b.tb.Send(m.Sender, err.Error(), tb.ModeMarkdown)
		return
	}

	s, err := b.subsApp.GetSubscription(params[0])
	if err != nil {
		log.Printf("failed to fetch subscription details, %s", err.Error())
		return
	}

	if s == nil {
		b.tb.Send(m.Sender, fmt.Sprintf("Cannot find subscription %s", params[0]))
		return
	}

	b.tb.Send(m.Sender, fmt.Sprintf("Subscription Details\n ```\n%s```", s.ToString()), tb.ModeMarkdown)
}

func (b Bot) subscribeForMovementCMD(m *tb.Message) {
	params, err := commands["subscribe"].ValidateParameters(m.Payload)
	if err != nil {
		b.tb.Send(m.Sender, err.Error(), tb.ModeMarkdown)
		return
	}

	for _, account := range params[1:] {
		if err := b.subsApp.Subscribe(
			m.Sender.Recipient(),
			params[0],
			account,
		); err != nil {
			log.Printf("failed to subscribe for movement, %s", err.Error())
			return
		}
	}

	b.tb.Send(
		m.Sender,
		fmt.Sprintf("subscribed to: `%s` accounts ```\n%+v\n``` for movement changes", params[0], params[1:]),
		tb.ModeMarkdown,
	)
}

func (b Bot) unsubscribeCMD(m *tb.Message) {
	params, err := commands["unsubscribe"].ValidateParameters(m.Payload)
	if err != nil {
		b.tb.Send(m.Sender, err.Error(), tb.ModeMarkdown)
		return
	}

	if err := b.subsApp.Unsubscribe(params[0]); err != nil {
		b.tb.Send(m.Sender, fmt.Sprintf("failed to unsubscribe, %s", err.Error()))
	}
}

func (b Bot) unsubscribeAllCMD(m *tb.Message) {
	if err := b.subsApp.UnsubscribeAllForUser(m.Sender.Recipient()); err != nil {
		b.tb.Send(m.Sender, fmt.Sprintf("failed to unsubscribe all, %s", err.Error()))
	}
}

func (b Bot) mySubscriptionsCMD(m *tb.Message) {
	subs, err := b.subsApp.GetSubscriptionsForUser(m.Sender.Recipient())
	if err != nil {
		log.Printf("failed to fetch subscriptions, %s", err.Error())
		return
	}

	msg := ""
	if len(subs) < 1 {
		msg = "I don't have any subscriptions"
	} else {
		for _, s := range subs {
			msg += fmt.Sprintf("%s\n\n", s.ToString())
		}
		msg = fmt.Sprintf("My Subscriptions: \n\n`%s`", msg)
	}

	b.tb.Send(m.Sender, msg, tb.ModeMarkdown)
}

func (b Bot) availableAssetsCMD(m *tb.Message) {
	s := "Avialable Assets\n\n\n```\n Bitcoin[btc], Ethereum[eth]```"
	b.tb.Send(m.Sender, s, tb.ModeMarkdown)
}

func (b Bot) availableCommandsCMD(m *tb.Message) {
	s := "Available Commands\n\n\n```\n"
	for cmdName, cmd := range commands {
		s += fmt.Sprintf("%s :\r\r %s\n\n", cmdName, cmd.Usage)
	}
	s += "```"
	b.tb.Send(m.Sender, s, tb.ModeMarkdown)
}

func (b Bot) helpCMD(m *tb.Message) {
	s := "Help\n\n\n```\n Shows this message```"
	b.tb.Send(m.Sender, s, tb.ModeMarkdown)
}

func (b Bot) defaultHandler(m *tb.Message) {
	s := fmt.Sprintf("Unknown command: ```\n %s```", m.Text)
	log.Print(s)
	b.tb.Send(m.Sender, s, tb.ModeMarkdown)
}
