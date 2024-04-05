package config

import (
	"flag"
	"github.com/caarlos0/env/v10"
	"log"
	"time"
)

type Config struct {
	RunAddr         string        `env:"RUN_ADDRESS"`
	LogLevel        string        `env:"LOG_LEVEL"`
	AccrualSysAddr  string        `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DBURI           string        `env:"DATABASE_URI"`
	JWTSecret       string        `env:"JWT_SECRET"`
	RequestInterval time.Duration `env:"REQUEST_INTERVAL"`
}

func New() *Config {

	flagCfg := &Config{}
	var reqInterval int
	flag.StringVar(&flagCfg.RunAddr, "a", ":8081", "run address")
	flag.StringVar(&flagCfg.LogLevel, "l", "debug", "log level")
	flag.StringVar(&flagCfg.AccrualSysAddr, "r", "http://localhost:8080", "accrual system address")
	flag.StringVar(&flagCfg.DBURI, "d", "postgres://postgres:postgres@localhost:5432/loyalty-system?sslmode=disable", "database uri")
	flag.StringVar(&flagCfg.JWTSecret, "s", "super-secret", "jwt secret")
	flag.IntVar(&reqInterval, "w", 5, "accrual request interval")
	flag.Parse()

	envCfg := &Config{}
	if err := env.Parse(envCfg); err != nil {
		log.Fatalf("Failed to parse env config: %v", err)
	}

	cfg := &Config{}
	cfg.RunAddr = envCfg.RunAddr
	cfg.LogLevel = envCfg.LogLevel
	cfg.DBURI = envCfg.DBURI
	cfg.AccrualSysAddr = envCfg.AccrualSysAddr
	cfg.JWTSecret = envCfg.JWTSecret
	if cfg.RunAddr == "" {
		cfg.RunAddr = flagCfg.RunAddr
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = flagCfg.LogLevel
	}
	if cfg.DBURI == "" {
		cfg.DBURI = flagCfg.DBURI
	}
	if cfg.AccrualSysAddr == "" {
		cfg.AccrualSysAddr = flagCfg.AccrualSysAddr
	}
	if cfg.JWTSecret == "" {
		cfg.JWTSecret = flagCfg.JWTSecret
	}
	if cfg.RequestInterval == 0 {
		cfg.RequestInterval = time.Duration(reqInterval)
	}

	return cfg
}
