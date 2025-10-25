package mocks

import (
	"context"
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/user/dao"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, user dao.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) FindByUsername(ctx context.Context, username string) (*dao.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*dao.User), args.Error(1)
}
