package handlers

import (
	"context"
	"github.com/vindosVP/loyalty-system/internal/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Storage
type Storage interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error)
	GetUsersOrders(ctx context.Context, userID int) ([]*models.Order, error)
	GetUsersCurrentBalance(ctx context.Context, userID int) (float64, error)
	GetUsersWithdrawnBalance(ctx context.Context, userID int) (float64, error)
	GetUsersWithdrawals(ctx context.Context, userID int) ([]*models.Order, error)
}
