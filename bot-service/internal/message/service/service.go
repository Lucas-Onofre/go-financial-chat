package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/Lucas-Onofre/financial-chat/bot-service/internal/broker"
	"github.com/Lucas-Onofre/financial-chat/bot-service/internal/message/dto"
	"github.com/Lucas-Onofre/financial-chat/bot-service/internal/shared"
	"github.com/Lucas-Onofre/financial-chat/bot-service/internal/stooqprovider"
)

var (
	ReasonExternalServiceFailure = "external service failure"
	ReasonInternalError          = "internal error, please try again later"
)

type Service struct {
	stooqClient    *stooqprovider.Client
	brokerProducer broker.Producer
}

func New(stooqClient *stooqprovider.Client, brokerProducer broker.Producer) *Service {
	return &Service{
		stooqClient:    stooqClient,
		brokerProducer: brokerProducer,
	}
}

func (s *Service) Process(_ context.Context, msg dto.CommandMessage) error {
	rawCsv, err := s.stooqClient.GetMarketData(msg.Command.GetValue())
	if err != nil {
		sendFailureMessage(s.brokerProducer, msg.UserID, msg.RoomID, ReasonExternalServiceFailure)
		return err
	}

	formattedMessage, err := getMessageFromCSV(rawCsv)
	if err != nil {
		sendFailureMessage(s.brokerProducer, msg.UserID, msg.RoomID, ReasonInternalError)
		return err
	}

	response := dto.ResponseMessage{
		UserID:  msg.UserID,
		RoomID:  msg.RoomID,
		Message: formattedMessage,
	}

	respBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("failed to marshal response: %v", err)
		return err
	}

	return s.brokerProducer.Publish(shared.BrokerChatResponsesQueueName, string(respBytes))
}

func getMessageFromCSV(csvData string) (string, error) {
	reader := csv.NewReader(strings.NewReader(csvData))
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	if len(records) < 2 || len(records[1]) < 7 {
		return "", errors.New("unexpected CSV format")
	}

	symbol := records[1][0]
	closePrice := records[1][6]

	if closePrice == "N/D" {
		return "", errors.New("quote not available")
	}

	return symbol + " quote is $" + closePrice + "per share", nil
}

func sendFailureMessage(broker broker.Producer, userID, roomID, reason string) {
	response := dto.ResponseMessage{
		UserID:  userID,
		RoomID:  roomID,
		Message: "Failed to retrieve stock data: " + reason,
	}

	respBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("failed to marshal failure response: %v", err)
		return
	}

	if err := broker.Publish(shared.BrokerChatResponsesQueueName, string(respBytes)); err != nil {
		log.Printf("failed to publish failure response: %v", err)
	}
}
