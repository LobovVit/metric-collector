package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Host           string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
	LogLevel       string `env:"LOG_LEVEL"`
	ReportFormat   string `env:"REPORT_FORMAT"`
}

func GetConfig() (*Config, error) {
	config := &Config{}
	err := env.Parse(config)
	if err != nil {
		return nil, fmt.Errorf("env parse failed: %w", err)
	}

	host := flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	reportInterval := flag.Int64("r", 10, "частота отправки метрик на сервер")
	pollInterval := flag.Int64("p", 2, "частота опроса метрик из пакета runtime")
	logLevel := flag.String("l", "info", "log level")
	reportFormat := flag.String("f", "json", "формат передачи метрик json/text")
	flag.Parse()

	if config.Host == "" {
		config.Host = *host
	}
	if config.Host != "" {
		config.Host = "http://" + config.Host + "/update"
	}
	if config.ReportInterval == 0 {
		config.ReportInterval = *reportInterval
	}
	if config.PollInterval == 0 {
		config.PollInterval = *pollInterval
	}
	if config.LogLevel == "" {
		config.LogLevel = *logLevel
	}
	if config.ReportFormat == "" {
		config.ReportFormat = *reportFormat
	}
	return config, nil
}
