package config

import (
	"flag"
	"github.com/caarlos0/env/v10"
	"log"
)

type Config struct {
	RunAddr        string `env:"RUN_ADDRESS"`
	LogLevel       string `env:"LOG_LEVEL"`
	AccrualSysAddr string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DbURI          string `env:"DATABASE_URI"`
	JWTSecret      string `env:"JWT_SECRET"`
}

func New() *Config {

	flagCfg := &Config{}
	flag.StringVar(&flagCfg.RunAddr, "a", ":8080", "run address")
	flag.StringVar(&flagCfg.LogLevel, "l", "debug", "log level")
	flag.StringVar(&flagCfg.AccrualSysAddr, "r", ":8081", "accrual system address")
	flag.StringVar(&flagCfg.DbURI, "d", "postgres://postgres:postgres@localhost:5432/loyalty-system?sslmode=disable", "database uri")
	flag.StringVar(&flagCfg.JWTSecret, "s", "super-secret", "jwt secret")
	flag.Parse()

	envCfg := &Config{}
	if err := env.Parse(envCfg); err != nil {
		log.Fatalf("Failed to parse env config: %v", err)
	}

	cfg := &Config{}
	cfg.RunAddr = envCfg.RunAddr
	cfg.LogLevel = envCfg.LogLevel
	cfg.DbURI = envCfg.DbURI
	cfg.AccrualSysAddr = envCfg.AccrualSysAddr
	cfg.JWTSecret = envCfg.JWTSecret
	if cfg.RunAddr == "" {
		cfg.RunAddr = flagCfg.RunAddr
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = flagCfg.LogLevel
	}
	if cfg.DbURI == "" {
		cfg.DbURI = flagCfg.DbURI
	}
	if cfg.AccrualSysAddr == "" {
		cfg.AccrualSysAddr = flagCfg.AccrualSysAddr
	}
	if cfg.JWTSecret == "" {
		cfg.JWTSecret = flagCfg.JWTSecret
	}

	return cfg
}
