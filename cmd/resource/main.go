package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/psychoplasma/crypto-balance-bot/application"
	"github.com/psychoplasma/crypto-balance-bot/infrastructure/persistence/mongodb"
)

var subsApp *application.SubscriptionApplication

func main() {
	subsRepo := mongodb.NewSubscriptionRepository()
	if err := subsRepo.Connect("mongodb://127.0.0.1:27017", "CryptoBalanceBot"); err != nil {
		panic(err)
	}
	defer subsRepo.Disconnect()
	subsApp = application.NewSubscriptionApplication(subsRepo)

	listenAndServe()
}

func listenAndServe() {
	log.Println("Starting API server on  localhost:1234")

	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/assets", GetAvailableAssets).Methods("GET")
	router.HandleFunc("/subscriptions/user/{userID}", GetSubscriptionsForUser).Methods("GET")
	router.HandleFunc("/subscriptions/{id}", GetSubscription).Methods("GET")

	log.Fatal(http.ListenAndServe("localhost:1234", router))
}
