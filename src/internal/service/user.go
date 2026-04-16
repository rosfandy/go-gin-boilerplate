package service

import (
	"context"

	"owner-api-proxy/internal/database/model"
)

type UserRepository interface {
	GetUsers(ctx context.Context) ([]model.Users, error)
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetUsers(ctx context.Context) ([]model.Users, error) {
	return s.userRepository.GetUsers(ctx)
}
