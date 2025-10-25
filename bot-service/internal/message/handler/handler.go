package handler

import (
	"context"
	"encoding/json"

	"github.com/Lucas-Onofre/financial-chat/bot-service/internal/message/dto"
	"github.com/Lucas-Onofre/financial-chat/bot-service/internal/message/service"
)

type Handler struct {
	service service.Service
}

func New(service service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Handle(ctx context.Context, message string) error {
	var commandMsg dto.CommandMessage
	if err := json.Unmarshal([]byte(message), &commandMsg); err != nil {
		return err
	}

	if validateErr := commandMsg.Validate(); validateErr != nil {
		return validateErr
	}

	return h.service.Process(ctx, commandMsg)
}
