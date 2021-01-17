package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/psychoplasma/crypto-balance-bot/application"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/persistence/mongodb"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/telegram"
	"gopkg.in/yaml.v2"
)

// Config represents configuration options for the observer
type Config struct {
	Token       string        `yaml:"token"`
	PollingTime time.Duration `yaml:"polling-time"`
	Database    struct {
		Name string `yaml:"name"`
		URI  string `yaml:"uri"`
	} `yaml:"database"`
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

func main() {
	c, err := readConfig("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	subsRepo, err := mongodb.NewSubscriptionRepository(c.Database.URI, c.Database.Name)
	if err != nil {
		panic(err)
	}

	o := NewMovementObserver(
		application.NewSubscriptionApplication(subsRepo),
		telegram.NewPublisher(c.Token, telegram.MovementFormatter),
	)

	sig := make(chan os.Signal)
	// Check for interrupt and kill signals so that we stop observer gracefully
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-sig:
			log.Println("interrupt received, exiting observer")
			o.Stop()
		}
	}()

	o.Start()
}
