package main

import (
	"context"
	"github.com/LobovVit/metric-collector/internal/agent/app"
	"github.com/LobovVit/metric-collector/internal/agent/config"
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
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGABRT)
	defer stop()
	agent := app.NewAgent(cfg, ctx)
	agent.RunAgent()
	return nil
}
