package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

type Config struct {
	Host           string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
}

func GetConfig() (*Config, error) {
	config := &Config{}
	err := env.Parse(config)
	if err != nil {
		return nil, errors.Wrap(err, "env parse failed")
	}

	host := flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	reportInterval := flag.Int64("r", 10, "частота отправки метрик на сервер")
	pollInterval := flag.Int64("p", 2, "частота опроса метрик из пакета runtime")
	flag.Parse()

	if config.Host == "" {
		config.Host = *host
	}
	if config.Host != "" {
		config.Host = "http://" + config.Host + "/update/"
	}
	if config.ReportInterval == 0 {
		config.ReportInterval = *reportInterval
	}
	if config.PollInterval == 0 {
		config.PollInterval = *pollInterval
	}
	return config, nil
}
