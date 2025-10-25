package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	brokermock "github.com/Lucas-Onofre/financial-chat/bot-service/internal/broker/mocks"
	mktdatamock "github.com/Lucas-Onofre/financial-chat/bot-service/internal/marketdataprovider/mocks"
	"github.com/Lucas-Onofre/financial-chat/bot-service/internal/message/dto"
	"github.com/Lucas-Onofre/financial-chat/bot-service/internal/shared"
)

var (
	validCSVResponse   = "Symbol,Date,Time,Open,High,Low,Close,Volume\nAAPL.US,2025-10-24,22:00:17,261.19,264.13,259.18,262.82,38253717"
	invalidCSVResponse = "Invalid,CSV,Data"
)

func TestService_Process(t *testing.T) {
	type args struct {
		ctx context.Context
		msg dto.CommandMessage
	}

	type want struct {
		error error
	}

	tests := []struct {
		name  string
		args  args
		setup func(marketDataClient *mktdatamock.MockMarketDataProvider, brokerProducer *brokermock.MockBroker)
		want  want
	}{
		{
			name: "Given valid command message, When Process is called, Then it should publish the correct response message",
			args: args{
				ctx: context.Background(),
				msg: dto.CommandMessage{
					UserID:  "user1",
					RoomID:  "room1",
					Command: dto.Command("/stock=AAPL"),
				},
			},
			setup: func(mktdataClient *mktdatamock.MockMarketDataProvider, brokerProducer *brokermock.MockBroker) {
				mktdataClient.On("GetMarketData", "AAPL").Return("Symbol,Date,Time,Open,High,Low,Close,Volume\nAAPL.US,2025-10-24,22:00:17,261.19,264.13,259.18,262.82,38253717", nil)
				brokerProducer.On("Publish", shared.BrokerChatResponsesQueueName, mock.Anything).Return(nil)
			},
			want: want{
				error: nil,
			},
		},
		{
			name: "Given market data provider failure, When Process is called, Then it should return an error and send a failure message",
			args: args{
				ctx: context.Background(),
				msg: dto.CommandMessage{
					UserID:  "user1",
					RoomID:  "room1",
					Command: dto.Command("/stock=AAPL"),
				},
			},
			setup: func(mktdataClient *mktdatamock.MockMarketDataProvider, brokerProducer *brokermock.MockBroker) {
				mktdataClient.On("GetMarketData", "AAPL").Return("", errors.New("error fetching market data"))
				brokerProducer.On("Publish", shared.BrokerChatResponsesQueueName, mock.Anything).Return(nil)
			},
			want: want{
				error: errors.New("error fetching market data"),
			},
		},
		{
			name: "Given invalid CSV data from market data provider, When Process is called, Then it should return an error and send a failure message",
			args: args{
				ctx: context.Background(),
				msg: dto.CommandMessage{
					UserID:  "user1",
					RoomID:  "room1",
					Command: dto.Command("/stock=AAPL"),
				},
			},
			setup: func(mktdataClient *mktdatamock.MockMarketDataProvider, brokerProducer *brokermock.MockBroker) {
				mktdataClient.On("GetMarketData", "AAPL").Return(invalidCSVResponse, nil)
				brokerProducer.On("Publish", shared.BrokerChatResponsesQueueName, mock.Anything).Return(nil)
			},
			want: want{
				error: errors.New("unexpected CSV format"),
			},
		},
		{
			name: "Given valid MarkedData but Publish fails, When Process is called, Then it should return an error",
			args: args{
				ctx: context.Background(),
				msg: dto.CommandMessage{
					UserID:  "user1",
					RoomID:  "room1",
					Command: dto.Command("/stock=AAPL"),
				},
			},
			setup: func(mktdataClient *mktdatamock.MockMarketDataProvider, brokerProducer *brokermock.MockBroker) {
				mktdataClient.On("GetMarketData", "AAPL").Return(validCSVResponse, nil)
				brokerProducer.On("Publish", shared.BrokerChatResponsesQueueName, mock.Anything).Return(errors.New("publish error"))
			},
			want: want{
				error: errors.New("publish error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mktdataClient := new(mktdatamock.MockMarketDataProvider)
			brokerProducer := new(brokermock.MockBroker)

			if tt.setup != nil {
				tt.setup(mktdataClient, brokerProducer)
			}

			service := New(mktdataClient, brokerProducer)
			err := service.Process(tt.args.ctx, tt.args.msg)

			assert.Equal(t, tt.want.error, err)
		})
	}
}

func TestService_sendFailureMessage(t *testing.T) {
	type args struct {
		userID string
		roomID string
		reason string
		broker *brokermock.MockBroker
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Given valid inputs, When sendFailureMessage is called, Then it should publish the failure message",
			args: args{
				userID: "user1",
				roomID: "room1",
				reason: "TestReason",
				broker: func() *brokermock.MockBroker {
					brokerProducer := new(brokermock.MockBroker)
					brokerProducer.On("Publish", shared.BrokerChatResponsesQueueName, mock.Anything).Return(nil)
					return brokerProducer
				}(),
			},
			wantErr: nil,
		},
		{
			name: "Given Publish fails, When sendFailureMessage is called, Then it should return an error",
			args: args{
				userID: "user1",
				roomID: "room1",
				reason: "TestReason",
				broker: func() *brokermock.MockBroker {
					brokerProducer := new(brokermock.MockBroker)
					brokerProducer.On("Publish", shared.BrokerChatResponsesQueueName, mock.Anything).Return(errors.New("publish error"))
					return brokerProducer
				}(),
			},
			wantErr: errors.New("publish error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sendFailureMessage(tt.args.broker, tt.args.userID, tt.args.roomID, tt.args.reason)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestService_getMessageFromCSV(t *testing.T) {
	type args struct {
		csvData string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "Given valid CSV data, When getMessageFromCSV is called, Then it should return the formatted message",
			args: args{
				csvData: validCSVResponse,
			},
			want:    "AAPL.US quote is $262.82 per share",
			wantErr: nil,
		},
		{
			name: "Given invalid CSV data, When getMessageFromCSV is called, Then it should return an error",
			args: args{
				csvData: invalidCSVResponse,
			},
			want:    "",
			wantErr: errors.New("unexpected CSV format"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getMessageFromCSV(tt.args.csvData)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
