package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/LobovVit/metric-collector/internal/agent/app"
	"github.com/LobovVit/metric-collector/internal/agent/config"
	"github.com/LobovVit/metric-collector/pkg/logger"
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
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		return fmt.Errorf("log initialize: %w", err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGABRT)
	defer cancel()
	agent := app.New(cfg)
	return agent.Run(ctx)
}
