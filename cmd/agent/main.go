// Package main - contains the main function that runs the agent application
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

// main - function that is executed at startup and is the entry point to the program
func main() {
	if err := run(context.Background()); err != nil {
		panic(err)
	}
}

// run - function starts the application, initializes the config and logger, creates an instance of the agent and launches it
func run(ctx context.Context) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		return fmt.Errorf("log initialize: %w", err)
	}
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGABRT)
	defer cancel()
	agent := app.New(cfg)
	return agent.Run(ctx)
}
