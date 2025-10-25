package main

import (
	"log"
	"net/http"

	"github.com/Lucas-Onofre/financial-chat/bot-service/internal/broker"
)

func main() {
	rb, err := broker.NewRabbitMQBroker("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer rb.Close()

	//stooqClient := stooqprovider.New()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
