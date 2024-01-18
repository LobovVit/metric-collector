package main

import (
	"github.com/LobovVit/metric-collector/internal/server/config"
	"github.com/LobovVit/metric-collector/internal/server/server"
	"github.com/pkg/errors"
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
	app := server.GetApp(cfg.Host)
	return app.RouterRun()
}
