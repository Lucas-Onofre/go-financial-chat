package handler

import (
	"encoding/json"
	customerrors "github.com/Lucas-Onofre/financial-chat/chat-service/internal/shared/errors"
	dto2 "github.com/Lucas-Onofre/financial-chat/chat-service/internal/user/dto"
	usersrv "github.com/Lucas-Onofre/financial-chat/chat-service/internal/user/service"
	"net/http"
)

var (
	ErrInvalidRequestBody = customerrors.AppError{
		Code:    "BAD_REQUEST",
		Message: "invalid request body",
		Status:  http.StatusBadRequest,
	}
)

type Handler struct {
	service usersrv.Service
}

func New(service usersrv.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (a *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input dto2.RegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrInvalidRequestBody)
		return
	}

	if err := a.service.Register(ctx, input); err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input dto2.LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrInvalidRequestBody)
	}

	token, err := a.service.Login(ctx, input)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto2.TokenDTO{TokenString: token})
}
