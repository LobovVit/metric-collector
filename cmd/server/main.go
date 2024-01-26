package main

import (
	"context"
	"fmt"
	"github.com/LobovVit/metric-collector/internal/server/config"
	"github.com/LobovVit/metric-collector/internal/server/logger"
	"github.com/LobovVit/metric-collector/internal/server/server"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("get config failed: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGABRT)
	defer cancel()

	if err := logger.Initialize(cfg.LogLevel); err != nil {
		return fmt.Errorf("log initialize failed: %w", err)
	}

	app := server.GetApp(cfg)
	return app.RouterRun(ctx)
}
