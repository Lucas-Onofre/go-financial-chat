package repository

import (
	"context"
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/user/dao"
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/user/repository/mocks"
	"testing"

	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Create(t *testing.T) {
	type args struct {
		context context.Context
		user    dao.User
	}

	tests := []struct {
		name    string
		args    args
		setup   func(repo *mocks.MockDB)
		wantErr bool
	}{
		{
			name: "Given valid user, When Create is called, Then no error is returned",
			args: args{
				context: context.Background(),
				user: dao.User{
					Username: "testuser",
					Password: "hashedpassword",
				},
			},
			setup: func(repo *mocks.MockDB) {
				repo.On("Create", mock.AnythingOfType("*dao.User")).Return(&gorm.DB{Error: nil})
			},
			wantErr: false,
		},
		{
			name: "Given DB error, When Create is called, Then error is returned",
			args: args{
				context: context.Background(),
				user: dao.User{
					Username: "testuser",
					Password: "hashedpassword",
				},
			},
			setup: func(repo *mocks.MockDB) {
				repo.On("Create", mock.AnythingOfType("*dao.User")).Return(&gorm.DB{Error: gorm.ErrInvalidData})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(mocks.MockDB)
			tt.setup(mockDB)

			repo := NewRepository(mockDB)
			err := repo.Create(tt.args.context, tt.args.user)

			assert.Equal(t, tt.wantErr, err != nil)
			mockDB.AssertExpectations(t)
		})
	}
}
