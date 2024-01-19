package main

import (
	"fmt"
	"github.com/LobovVit/metric-collector/internal/server/config"
	"github.com/LobovVit/metric-collector/internal/server/server"
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
	app := server.GetApp(cfg.Host)
	return app.RouterRun()
}
