package repos

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vindosVP/loyalty-system/internal/models"
)

type OrdersRepo struct {
	pool *pgxpool.Pool
}

func NewOrdersRepo(pool *pgxpool.Pool) *OrdersRepo {
	return &OrdersRepo{pool: pool}
}

func (or *OrdersRepo) Create(ctx context.Context, order *models.Order) (*models.Order, error) {
	query := "insert into orders (id, user_id, status, sum, uploaded_at) values ($1, $2, $3, $4, $5)"
	_, err := or.pool.Exec(ctx, query, order.ID, order.UserID, order.Status, order.Sum, order.UploadedAt)
	if err != nil {
		return nil, fmt.Errorf("error or.pool.Exec: %w", err)
	}
	resOrder, err := or.GetByID(ctx, order.ID)
	if err != nil {
		return nil, fmt.Errorf("or.GetByID: %w", err)
	}
	return resOrder, nil
}

func (or *OrdersRepo) GetByID(ctx context.Context, id int) (*models.Order, error) {
	query := "select id, user_id, status, sum, uploaded_at from orders where id = $1"
	row := or.pool.QueryRow(ctx, query, id)
	order := &models.Order{}
	err := row.Scan(&order.ID, &order.UserID, &order.Status, &order.Sum, &order.UploadedAt)
	if err != nil {
		return nil, fmt.Errorf("row.Scan: %w", err)
	}
	return order, nil
}

func (or *OrdersRepo) Exists(ctx context.Context, id int) (bool, error) {
	query := "select exists(select 1 from orders where id = $1)"
	row := or.pool.QueryRow(ctx, query, id)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("row.Scan: %w", err)
	}
	return exists, nil
}

func (or *OrdersRepo) GetUsersOrders(ctx context.Context, userID int) ([]*models.Order, error) {
	query := "select id, user_id, status, sum, uploaded_at from orders where user_id = $1 order by uploaded_at"
	orders := make([]*models.Order, 0)
	rows, err := or.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("or.pool.Query: %w", err)
	}
	for rows.Next() {
		order := &models.Order{}
		err := rows.Scan(&order.ID, &order.UserID, &order.Status, &order.Sum, &order.UploadedAt)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (or *OrdersRepo) GetUsersCurrentBalance(ctx context.Context, userID int) (float64, error) {
	query := "select sum(sum) from orders where user_id = $1"
	row := or.pool.QueryRow(ctx, query, userID)
	var balance sql.NullFloat64
	err := row.Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("row.Scan: %w", err)
	}
	if balance.Valid {
		return balance.Float64, nil
	}
	return 0, nil
}

func (or *OrdersRepo) GetUsersWithdrawnBalance(ctx context.Context, userID int) (float64, error) {
	query := "select -sum(sum) from orders where user_id = $1 and sum < 0"
	row := or.pool.QueryRow(ctx, query, userID)
	var balance sql.NullFloat64
	err := row.Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("row.Scan: %w", err)
	}
	if balance.Valid {
		return balance.Float64, nil
	}
	return 0, nil
}
