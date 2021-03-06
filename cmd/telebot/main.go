package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/psychoplasma/crypto-balance-bot/application"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/persistence/mongodb"
	"gopkg.in/yaml.v2"
)

// Config is a configuration for telegram bot
type Config struct {
	Token       string        `yaml:"token"`
	PollingTime time.Duration `yaml:"polling-time"`
}

func main() {
	c, err := readConfig("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// var subsRepo = inmemory.NewSubscriptionRepository()
	subsRepo, err := mongodb.NewSubscriptionRepository("mongodb://127.0.0.1:27017", "CryptoBalanceBot")
	if err != nil {
		panic(err)
	}
	var subsAppService = application.NewSubscriptionApplication(subsRepo)

	b := NewBot(c, subsAppService)
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
