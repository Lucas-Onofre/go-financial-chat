package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Lucas-Onofre/financial-chat/bot-service/internal/broker"
	"github.com/Lucas-Onofre/financial-chat/bot-service/internal/marketdataprovider"
	"github.com/Lucas-Onofre/financial-chat/bot-service/internal/message/handler"
	"github.com/Lucas-Onofre/financial-chat/bot-service/internal/message/service"
	"github.com/Lucas-Onofre/financial-chat/bot-service/internal/shared"
)

func main() {
	// TODO configure rabbitMQ via env variables
	rb, err := broker.NewRabbitMQBroker("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer rb.Close()

	stooqClient := marketdataprovider.New()
	service := service.New(stooqClient, rb)
	handler := handler.New(service)

	if err := rb.Subscribe(shared.BrokerChatCommandsQueueName, func(message string) error {
		if err := handler.Handle(context.Background(), message); err != nil {
			log.Printf("failed to handle message: %v", err)
			return err
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
