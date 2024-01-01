package main

import (
	"github.com/LobovVit/metric-collector/internal/agent/config"
	"github.com/LobovVit/metric-collector/internal/agent/skeduller"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	cfg := config.GetConfig()
	skeduller.StartTimer(cfg.PollInterval, cfg.ReportInterval, "http://"+cfg.Host+"/update/")
	return nil
}
