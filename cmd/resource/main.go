package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/psychoplasma/crypto-balance-bot/application"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/services"
	"gopkg.in/yaml.v2"
)

var subsApp *application.SubscriptionApplication

// Config is a configuration for telegram bot
type Config struct {
	Resource struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"resource"`
	Database struct {
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
	subsApp = application.NewSubscriptionApplication(subsRepo)

	listenAndServe(c.Resource.Host, c.Resource.Port)
}

func listenAndServe(host string, port int) {
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Starting API server on  %s\n\r", addr)

	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/assets", GetAvailableAssets).Methods("GET")
	router.HandleFunc("/subscriptions/user/{userID}", GetSubscriptionsForUser).Methods("GET")
	router.HandleFunc("/subscriptions/{id}", GetSubscription).Methods("GET")

	log.Fatal(http.ListenAndServe(addr, router))
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
