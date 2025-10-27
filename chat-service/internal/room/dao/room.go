package dao

import "github.com/Lucas-Onofre/financial-chat/chat-service/internal/shared/entity"

type Room struct {
	entity.Entity
	Name string `json:"name"`
}
