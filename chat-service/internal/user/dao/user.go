package dao

import (
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/shared/entity"
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/user/dto"
	"time"

	"github.com/google/uuid"
)

type User struct {
	entity.Entity
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`
}

func (u User) Build(password string) User {
	return User{
		Entity: entity.Entity{
			ID:        uuid.NewString(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: u.Username,
		Password: password,
	}
}

func (u User) FromRegisterDTO(dto dto.RegisterDTO) User {
	return User{
		Username: dto.Username,
		Password: dto.Password,
	}
}

func (u User) FromLoginDTO(dto dto.LoginDTO) User {
	return User{
		Username: dto.Username,
		Password: dto.Password,
	}
}
