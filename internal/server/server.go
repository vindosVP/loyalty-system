package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	chim "github.com/go-chi/chi/v5/middleware"
	"github.com/vindosVP/loyalty-system/cmd/gophermart/config"
	"github.com/vindosVP/loyalty-system/internal/database"
	"github.com/vindosVP/loyalty-system/internal/handlers"
	"github.com/vindosVP/loyalty-system/internal/middleware"
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
	or := repos.NewOrdersRepo(pool)
	s := storage.New(ur, or)

	r := chi.NewRouter()
	r.Use(chim.Logger, chim.Compress(5))
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	r.Post("/api/user/register", handlers.Register(s, cfg.JWTSecret))
	r.Post("/api/user/login", handlers.Login(s, cfg.JWTSecret))
	r.Group(func(r chi.Router) {
		a := middleware.NewAuthenticator(cfg.JWTSecret)
		r.Use(a.WithAuth)
		r.Post("/api/user/orders", handlers.CreateOrder(s))
		r.Get("/api/user/orders", handlers.GetOrderList(s))
		r.Get("/api/user/balance", handlers.GetUsersBalance(s))
		r.Post("/api/user/balance/withdraw", handlers.WithdrawOrder(s))
		r.Get("/api/user/withdrawals", handlers.GetUsersWithdrawals(s))
	})

	logger.Log.Info("Server started", zap.String("Address", cfg.RunAddr))
	err = http.ListenAndServe(cfg.RunAddr, r)
	if err != nil {
		return fmt.Errorf("http.ListenAndServe: %w", err)
	}
	return nil
}
