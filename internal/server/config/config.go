package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Host string `env:"ADDRESS"`
}

func GetConfig() (*Config, error) {
	config := &Config{}
	err := env.Parse(config)
	if err != nil {
		return nil, err
	}
	host := flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	flag.Parse()

	if config.Host == "" {
		config.Host = *host
	}
	return config, nil
}
