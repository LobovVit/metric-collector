package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/LobovVit/metric-collector/internal/server/config"
	"github.com/LobovVit/metric-collector/internal/server/server"
	"github.com/LobovVit/metric-collector/pkg/logger"
)

// main - function that is executed at startup and is the entry point to the program
func main() {
	if err := run(context.Background()); err != nil {
		panic(err)
	}
}

// run - function starts the application, initializes the config and logger, creates an instance of the server and launches it
func run(ctx context.Context) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}
	if err = logger.Initialize(cfg.LogLevel); err != nil {
		return fmt.Errorf("log initialize: %w", err)
	}
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGABRT)
	defer cancel()
	app, err := server.New(ctx, cfg)
	if err != nil {
		return err
	}
	return app.Run(ctx)
}
