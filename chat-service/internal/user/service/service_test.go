package service

import (
	"context"
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/auth/jwt"
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/auth/jwt/utils"
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/user/dao"
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/user/dto"
	userrepomock "github.com/Lucas-Onofre/financial-chat/chat-service/internal/user/repository/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Login(t *testing.T) {
	validPass, _ := utils.HashPassword("password123")

	type args struct {
		ctx      context.Context
		loginDTO dto.LoginDTO
	}
	tests := []struct {
		name    string
		args    args
		setup   func(repo *userrepomock.MockRepository)
		wantErr bool
	}{
		{
			name: "Given valid credentials, When Login is called, Then no error is returned",
			args: args{
				ctx: context.Background(),
				loginDTO: dto.LoginDTO{
					Username: "testuser",
					Password: "password123",
				},
			},
			setup: func(repo *userrepomock.MockRepository) {
				repo.On("FindByUsername", mock.Anything, "testuser").Return(&dao.User{
					Username: "testuser",
					Password: validPass,
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "Given invalid password, When Login is called, Then error is returned",
			args: args{
				ctx: context.Background(),
				loginDTO: dto.LoginDTO{
					Username: "testuser",
					Password: validPass,
				},
			},
			setup: func(repo *userrepomock.MockRepository) {
				repo.On("FindByUsername", mock.Anything, "testuser").Return(&dao.User{
					Username: "testuser",
					Password: "invalid",
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "Given non-existing user, When Login is called, Then error is returned",
			args: args{
				ctx: context.Background(),
				loginDTO: dto.LoginDTO{
					Username: "nonexistinguser",
					Password: "pass",
				},
			},
			setup: func(repo *userrepomock.MockRepository) {
				repo.On("FindByUsername", mock.Anything, "nonexistinguser").Return(new(dao.User), nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(userrepomock.MockRepository)
			tt.setup(mockRepo)

			jwtService := jwt.NewJWTService("secret", time.Minute*5)
			service := New(mockRepo, jwtService)

			_, err := service.Login(tt.args.ctx, tt.args.loginDTO)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestService_Register(t *testing.T) {
	type args struct {
		ctx         context.Context
		registerDTO dto.RegisterDTO
	}
	tests := []struct {
		name    string
		args    args
		setup   func(repo *userrepomock.MockRepository)
		wantErr bool
	}{
		{
			name: "Given valid user, When Register is called, Then no error is returned",
			args: args{
				ctx: context.Background(),
				registerDTO: dto.RegisterDTO{
					Username: "newuser",
					Password: "password123",
				},
			},
			setup: func(repo *userrepomock.MockRepository) {
				repo.On("Create", mock.Anything, mock.MatchedBy(func(u dao.User) bool {
					return u.Username == "newuser" && len(u.Password) > 0
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Given repository error, When Register is called, Then error is returned",
			args: args{
				ctx: context.Background(),
				registerDTO: dto.RegisterDTO{
					Username: "newuser",
					Password: "password123",
				},
			},
			setup: func(repo *userrepomock.MockRepository) {
				repo.On("Create", mock.Anything, mock.MatchedBy(func(u dao.User) bool {
					return u.Username == "newuser" && len(u.Password) > 0
				})).Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(userrepomock.MockRepository)
			tt.setup(mockRepo)

			service := New(mockRepo, nil)

			err := service.Register(tt.args.ctx, tt.args.registerDTO)
			assert.Equal(t, tt.wantErr, err != nil)
			mockRepo.AssertExpectations(t)
		})
	}
}
