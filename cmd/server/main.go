package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/LobovVit/metric-collector/internal/server/config"
	"github.com/LobovVit/metric-collector/internal/server/server"
	"github.com/LobovVit/metric-collector/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}
	if err = logger.Initialize(cfg.LogLevel); err != nil {
		return fmt.Errorf("log initialize: %w", err)
	}
	logger.Log.Info("Config", zap.String("SigningKey", cfg.SigningKey),
		zap.String("DSN", cfg.DSN),
		zap.String("", cfg.LogLevel))
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGABRT)
	defer cancel()
	app, err := server.New(ctx, cfg)
	if err != nil {
		return err
	}
	return app.Run(ctx)
}
