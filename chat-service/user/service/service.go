package service

import (
	"context"
	"errors"
	"fmt"

	jwtport "github.com/Lucas-Onofre/financial-chat/chat-service/auth/jwt/port"
	"github.com/Lucas-Onofre/financial-chat/chat-service/auth/jwt/utils"
	customerrors "github.com/Lucas-Onofre/financial-chat/chat-service/shared/errors"
	"github.com/Lucas-Onofre/financial-chat/chat-service/user/dao"
	"github.com/Lucas-Onofre/financial-chat/chat-service/user/dto"
)

type Service struct {
	repo       any
	jwtService jwtport.TokenService
}

func New(repo any, jwtService jwtport.TokenService) *Service {
	return &Service{
		repo:       repo,
		jwtService: jwtService,
	}
}

func (s *Service) Register(ctx context.Context, userDTO dto.RegisterDTO) error {
	var user dao.User
	user = user.FromRegisterDTO(userDTO)

	hashedPassword, hashErr := utils.HashPassword(user.Password)
	if hashErr != nil {
		return hashErr
	}

	fmt.Println("Hashed password:", hashedPassword)
	return nil
	// TODO uncomment it after creating the repository
	//return s.repo.Create(ctx, user.Build(hashedPassword))
}

func (s *Service) Login(ctx context.Context, loginDTO dto.LoginDTO) (string, error) {
	var user dao.User
	user = user.FromLoginDTO(loginDTO)

	saved, err := s.repo.FindByUsername(ctx, user.Username)
	if err != nil || saved == nil {
		return "", customerrors.Wrap(customerrors.ErrUnauthorized, errors.New("error retrieving data"))
	}

	if !utils.CheckPasswordHash(loginDTO.Password, saved.Password) {
		return "", customerrors.Wrap(customerrors.ErrUnauthorized, errors.New("invalid credentials"))
	}

	return s.jwtService.GenerateToken(saved.ID)
}
