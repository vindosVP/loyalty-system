package storage

import (
	"context"
	"fmt"
	"github.com/vindosVP/loyalty-system/internal/models"
	"github.com/vindosVP/loyalty-system/internal/repos"
)

type Storage struct {
	userRepo  *repos.UserRepo
	orderRepo *repos.OrdersRepo
}

func New(ur *repos.UserRepo, or *repos.OrdersRepo) *Storage {
	return &Storage{userRepo: ur, orderRepo: or}
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

func (s *Storage) CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	orderExists, err := s.orderRepo.Exists(ctx, order.ID)
	if err != nil {
		return nil, fmt.Errorf("s.orderRepo.Exists: %w", err)
	}
	if orderExists {
		existingOrder, err := s.orderRepo.GetByID(ctx, order.ID)
		if err != nil {
			return nil, fmt.Errorf("s.orderRepo.GetByID: %w", err)
		}
		if existingOrder.UserID == order.UserID {
			return nil, ErrOrderAlreadyExists
		} else {
			return nil, ErrOrderCreatedByOtherUser
		}
	}
	newOrder, err := s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("s.orderRepo.Create: %w", err)
	}
	return newOrder, nil
}

func (s *Storage) GetUsersOrders(ctx context.Context, userID int) ([]*models.Order, error) {
	orders, err := s.orderRepo.GetUsersOrders(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("s.orderRepo.GetUsersOrders: %w", err)
	}
	return orders, nil
}

func (s *Storage) GetUsersCurrentBalance(ctx context.Context, userID int) (float64, error) {
	balance, err := s.orderRepo.GetUsersCurrentBalance(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("s.orderRepo.GetUsersBalance: %w", err)
	}
	return balance, nil
}

func (s *Storage) GetUsersWithdrawnBalance(ctx context.Context, userID int) (float64, error) {
	balance, err := s.orderRepo.GetUsersWithdrawnBalance(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("s.orderRepo.GetUsersWithdrawnBalance: %w", err)
	}
	return balance, nil
}

func (s *Storage) GetUsersWithdrawals(ctx context.Context, userID int) ([]*models.Order, error) {
	withdrawals, err := s.orderRepo.GetUsersWithdrawals(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("s.orderRepo.GetUsersWithdrawals: %w", err)
	}
	return withdrawals, nil
}
