package port

import (
	"context"

	"github.com/Lucas-Onofre/financial-chat/chat-service/user/dao"
)

type RepositoryPort interface {
	Create(ctx context.Context, user dao.User) error
	FindByUsername(ctx context.Context, username string) (*dao.User, error)
}
