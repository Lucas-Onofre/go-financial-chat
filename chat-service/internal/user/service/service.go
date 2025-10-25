package service

import (
	"context"
	"errors"
	jwtport "github.com/Lucas-Onofre/financial-chat/chat-service/internal/auth/jwt/port"
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/auth/jwt/utils"
	customerrors "github.com/Lucas-Onofre/financial-chat/chat-service/internal/shared/errors"
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/user/dao"
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/user/dto"
	userrepo "github.com/Lucas-Onofre/financial-chat/chat-service/internal/user/repository/port"
)

type Service struct {
	repo       userrepo.RepositoryPort
	jwtService jwtport.TokenService
}

func New(repo userrepo.RepositoryPort, jwtService jwtport.TokenService) *Service {
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

	return s.repo.Create(ctx, user.Build(hashedPassword))
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
