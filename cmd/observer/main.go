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
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/port/adapter/publisher/telegram"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/services"
	"gopkg.in/yaml.v2"
)

// Config represents configuration options for the observer
type Config struct {
	Telebot struct {
		Token string `yaml:"token"`
	} `yaml:"telebot"`
	Observer struct {
		Currency          string        `yaml:"currency"`
		BlockHeightMargin uint64        `yaml:"block-margin"`
		Interval          time.Duration `yaml:"interval"`
		Parallelism       int           `yaml:"parallelism"`
		ExitTimeout       time.Duration `yaml:"exit-timeout"`
	} `yaml:"observer"`
	Database struct {
		Type string `yaml:"type"`
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

	subsRepo := services.RepositoryServiceFactory[c.Database.Type]
	if subsRepo == nil {
		panic(fmt.Errorf("there is no repository implementation for the given database type(%s)", c.Database.Type))
	}

	if err := subsRepo.Connect(c.Database.URI, c.Database.Name); err != nil {
		panic(err)
	}
	defer subsRepo.Disconnect()

	o := NewMovementObserver(
		application.NewSubscriptionApplication(subsRepo),
		telegram.NewPublisher(c.Telebot.Token, telegram.MovementFormatter),
		c.Observer.Currency,
		&ObserverOptions{
			BlockHeightMargin: c.Observer.BlockHeightMargin,
			ObserveInterval:   c.Observer.Interval * time.Second,
			MaxParallelism:    c.Observer.Parallelism,
			ExitTimeout:       c.Observer.ExitTimeout * time.Second,
		},
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
