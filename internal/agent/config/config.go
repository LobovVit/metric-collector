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
	SigningKey     string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
	MaxCntInBatch  int    `env:"BATCH_LIMIT"`
}

func GetConfig() (*Config, error) {
	config := &Config{}
	err := env.Parse(config)
	if err != nil {
		return nil, fmt.Errorf("env parse: %w", err)
	}

	host := flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	reportInterval := flag.Int64("r", 10, "частота отправки метрик на сервер")
	pollInterval := flag.Int64("p", 2, "частота опроса метрик из пакета runtime")
	logLevel := flag.String("log", "info", "log level")
	reportFormat := flag.String("f", "batch", "формат передачи метрик json/text/batch")
	maxCntInBatch := flag.Int("m", 5, "максимальное количество метрик в батче")
	signingKey := flag.String("k", "", "ключ")
	rateLimit := flag.Int("l", 10, "максимальное кол-во одновременно исходящих запросов на сервер")
	flag.Parse()

	if config.ReportFormat == "" {
		config.ReportFormat = *reportFormat
	}
	if config.Host == "" {
		config.Host = *host
	}
	if config.Host != "" && config.ReportFormat != "batch" {
		config.Host = "http://" + config.Host + "/update/"
	}
	if config.Host != "" && config.ReportFormat == "batch" {
		config.Host = "http://" + config.Host + "/updates/"
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
	if config.SigningKey == "" {
		config.SigningKey = *signingKey
	}
	if config.RateLimit == 0 {
		config.RateLimit = *rateLimit
	}
	if config.MaxCntInBatch == 0 {
		config.MaxCntInBatch = *maxCntInBatch
	}

	return config, nil
}
