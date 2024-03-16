package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vindosVP/loyalty-system/cmd/gophermart/config"
	"github.com/vindosVP/loyalty-system/internal/database"
	"github.com/vindosVP/loyalty-system/internal/handlers"
	"github.com/vindosVP/loyalty-system/internal/repos"
	"github.com/vindosVP/loyalty-system/internal/storage"
	"github.com/vindosVP/loyalty-system/pkg/logger"
	"go.uber.org/zap"
	"net/http"
)

func Run(cfg *config.Config) error {
	ctx := context.Background()
	pool, err := database.New(ctx, cfg.DBURI)
	if err != nil {
		return fmt.Errorf("database.New: %w", err)
	}
	defer pool.Close()

	ur := repos.NewUserRepo(pool)
	s := storage.New(ur)

	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.Compress(5))
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	r.Post("/api/user/register", handlers.Register(s, cfg.JWTSecret))
	r.Post("/api/user/login", handlers.Login(s, cfg.JWTSecret))

	logger.Log.Info("Server started", zap.String("Address", cfg.RunAddr))
	err = http.ListenAndServe(cfg.RunAddr, r)
	if err != nil {
		return fmt.Errorf("http.ListenAndServe: %w", err)
	}
	return nil
}
