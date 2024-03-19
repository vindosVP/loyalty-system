package repos

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vindosVP/loyalty-system/internal/models"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (ur *UserRepo) Create(ctx context.Context, user *models.User) (*models.User, error) {
	query := "insert into users (login, encryptedPassword) values ($1, $2)"
	_, err := ur.pool.Exec(ctx, query, user.Login, user.EncryptedPwd)
	if err != nil {
		return nil, fmt.Errorf("ur.pool.Exec: %w", err)
	}
	resUser, err := ur.GetByLogin(ctx, user.Login)
	if err != nil {
		return nil, fmt.Errorf("ur.GetUserByLogin: %w", err)
	}
	return resUser, nil
}

func (ur *UserRepo) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	query := "select id, login, encryptedPassword from users where login = $1"
	row := ur.pool.QueryRow(ctx, query, login)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Login, &user.EncryptedPwd)
	if err != nil {
		return nil, fmt.Errorf("row.Scan: %w", err)
	}
	return user, nil
}

func (ur *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := "select id, login, encryptedPassword from users where id = $1"
	row := ur.pool.QueryRow(ctx, query, id)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Login, &user.EncryptedPwd)
	if err != nil {
		return nil, fmt.Errorf("row.Scan: %w", err)
	}
	return user, nil
}

func (ur *UserRepo) Exists(ctx context.Context, login string) (bool, error) {
	query := "select exists(select 1 from users where login = $1)"
	row := ur.pool.QueryRow(ctx, query, login)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("row.Scan: %w", err)
	}
	return exists, nil
}
