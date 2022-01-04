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
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/services"
	"gopkg.in/yaml.v2"
)

// Config is a configuration for telegram bot
type Config struct {
	Token       string        `yaml:"token"`
	PollingTime time.Duration `yaml:"polling-time"`
	Database    struct {
		Type string `yaml:"type"`
		Name string `yaml:"name"`
		URI  string `yaml:"uri"`
	} `yaml:"database"`
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
	var subsAppService = application.NewSubscriptionApplication(subsRepo)

	b := NewBot(c, subsAppService)

	sig := make(chan os.Signal)
	// Check for interrupt and kill signals so that we stop observer gracefully
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-sig:
			log.Println("interrupt received, exiting telebot")
			b.Stop()
		}
	}()

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
