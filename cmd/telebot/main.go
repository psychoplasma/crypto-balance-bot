package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/psychoplasma/crypto-balance-bot/application"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/notification"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/persistence/inmemory"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Token       string        `yaml:"token"`
	PollingTime time.Duration `yaml:"polling-time"`
}

var subsRepo = inmemory.NewSubscriptionReposititory()
var subsAppService = application.NewSubscriptionApplication(subsRepo)
var currencyAppService = application.NewCurrencyService()
var movementPublisher = adapter.NewTelegramPublisher(
	"873977886:AAEJetV4LiotkaqDo3NGOrZXQ2BWEA2U8ts", notification.MovementFormatter{})

func main() {
	c, err := readConfig("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	o := application.NewObserver(subsRepo)
	o.RegisterPublisher(adapter.NewTelegramPublisher(c.Token, notification.MovementFormatter{}))
	go o.Observe()

	b := NewBot(c, subsAppService, currencyAppService)
	b.Start()
}

func readConfig(path string) (*Config, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, fmt.Errorf("configuration file does not exist")
	}

	c := &Config{}
	err = yaml.Unmarshal(f, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
