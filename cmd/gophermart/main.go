package main

import (
	"github.com/vindosVP/loyalty-system/cmd/gophermart/config"
	"github.com/vindosVP/loyalty-system/internal/server"
	"github.com/vindosVP/loyalty-system/pkg/logger"
	"go.uber.org/zap"
	"log"
)

func main() {
	log.Print("Starting loyalty system")
	cfg := config.New()
	err := logger.Initialize(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	logger.Log.Info("Configuration loaded",
		zap.String("Run address", cfg.RunAddr),
		zap.String("Log level", cfg.LogLevel),
		zap.String("Database URI", cfg.DbURI),
		zap.String("Accrual system address", cfg.AccrualSysAddr),
	)
	err = server.Run(cfg)
	if err != nil {
		logger.Log.Fatal("Failed to run server", zap.Error(err))
	}
}
