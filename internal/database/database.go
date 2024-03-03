package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vindosVP/loyalty-system/pkg/logger"
)

func New(ctx context.Context, dbURI string) (*pgxpool.Pool, error) {
	logger.Log.Info("Connecting to database")
	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}
	logger.Log.Info("Connected successfully")
	logger.Log.Info("Creating tables")
	err = createTables(ctx, pool)
	if err != nil {
		return nil, fmt.Errorf("createTables: %w", err)
	}
	logger.Log.Info("Tables created successfully")
	return pool, nil
}

func createTables(ctx context.Context, pool *pgxpool.Pool) error {
	query := "CREATE TABLE IF NOT EXISTS users (id SERIAL NOT NULL PRIMARY KEY, login TEXT NOT NULL, encryptedPassword TEXT NOT NULL);"
	_, err := pool.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}
