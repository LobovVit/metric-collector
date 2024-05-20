// Package config - included struct and init function fow work with app configuration
package config

import (
	"flag"
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
)

// Config determines the basic parameters of the agent's operation
type Config struct {
	Host           string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
	LogLevel       string `env:"LOG_LEVEL"`
	ReportFormat   string `env:"REPORT_FORMAT"`
	SigningKey     string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
	MaxCntInBatch  int    `env:"BATCH_LIMIT"`
	CryptoKey      string `env:"CRYPTO_KEY"`
	ConfigPath     string `env:"CONFIG"`
}

// GetConfig - method creates a new configuration and sets values from environment variables and command line flags
func GetConfig() (*Config, error) {
	config := &Config{}
	err := env.Parse(config)
	if err != nil {
		return nil, fmt.Errorf("env parse: %w", err)
	}

	host := flag.String("a", "", "адрес эндпоинта HTTP-сервера")                  //localhost:8080
	reportInterval := flag.Int64("r", 10, "частота отправки метрик на сервер")    //10
	pollInterval := flag.Int64("p", 2, "частота опроса метрик из пакета runtime") //2
	logLevel := flag.String("log", "info", "log level")
	reportFormat := flag.String("f", "batch", "формат передачи метрик json/text/batch")
	maxCntInBatch := flag.Int("m", 5, "максимальное количество метрик в батче")
	signingKey := flag.String("k", "", "ключ")
	rateLimit := flag.Int("l", 10, "максимальное кол-во одновременно исходящих запросов на сервер")
	cryptoKey := flag.String("crypto-key", "", "путь до файла с публичным ключом") //public.pem
	configPath1 := flag.String("config", "", "файл с JSON конфигом")
	configPath2 := flag.String("c", "", "файл с JSON конфигом")
	flag.Parse()

	if config.ConfigPath == "" {
		if *configPath1 != "" {
			config.ConfigPath = *configPath1
		}
		if *configPath2 != "" {
			config.ConfigPath = *configPath2
		}
	}

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
	if config.CryptoKey == "" {
		config.CryptoKey = *cryptoKey
	}

	config, err = parseJSONConfig(*config)
	if err != nil {
		log.Printf("parseJSONConfig: %v", err)
	}
	return config, nil
}
