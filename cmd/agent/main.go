package main

import (
	"context"
	"github.com/LobovVit/metric-collector/internal/agent/app"
	"github.com/LobovVit/metric-collector/internal/agent/config"
	"github.com/pkg/errors"
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
		return errors.Wrap(err, "get config failed")
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGABRT)
	defer cancel()
	agent := app.NewAgent(cfg)
	agent.RunAgent(ctx)
	return nil
}
