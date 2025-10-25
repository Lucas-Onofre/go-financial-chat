package port

import (
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/auth/jwt/model"
)

type TokenService interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(token string) (*model.CustomClaims, error)
}
