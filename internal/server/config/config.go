package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Host     string `env:"ADDRESS"`
	LogLevel string `env:"LOG_LEVEL"`
}

func GetConfig() (*Config, error) {
	config := &Config{}
	err := env.Parse(config)
	if err != nil {
		return nil, fmt.Errorf("env parse failed: %w", err)
	}
	host := flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	logLevel := flag.String("l", "info", "log level")
	flag.Parse()

	if config.Host == "" {
		config.Host = *host
	}
	if config.LogLevel == "" {
		config.LogLevel = *logLevel
	}

	return config, nil
}
