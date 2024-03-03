package storage

import (
	"context"
	"fmt"
	"github.com/vindosVP/loyalty-system/internal/models"
	"github.com/vindosVP/loyalty-system/internal/repos"
)

type Storage struct {
	userRepo *repos.UserRepo
}

func New(ur *repos.UserRepo) *Storage {
	return &Storage{userRepo: ur}
}

func (s *Storage) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	userExists, err := s.userRepo.Exists(ctx, user.Login)
	if err != nil {
		return nil, fmt.Errorf("s.userRepo.Exists: %w", err)
	}
	if userExists {
		return nil, ErrUserAlreadyExists
	}
	newUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("s.userRepo.Create: %w", err)
	}
	return newUser, nil
}

func (s *Storage) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	userExists, err := s.userRepo.Exists(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("s.userRepo.Exists: %w", err)
	}
	if !userExists {
		return nil, ErrUserNotFound
	}
	user, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("s.userRepo.GetByLogin: %w", err)
	}
	return user, nil
}
