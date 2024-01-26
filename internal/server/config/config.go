package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"os"
)

type Config struct {
	Host            string `env:"ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}

func GetConfig() (*Config, error) {
	config := &Config{}
	err := env.Parse(config)
	if err != nil {
		return nil, fmt.Errorf("env parse failed: %w", err)
	}

	host := flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	logLevel := flag.String("l", "info", "log level")
	storeInterval := flag.Int("i", 30, "интервал сохранения на диск")
	fileStoragePath := flag.String("f", "/tmp/metrics-db.json", "файл для сохранения на диск")
	restore := flag.Bool("r", true, "загружать при старте данные из файла")
	flag.Parse()

	if config.Host == "" {
		config.Host = *host
	}

	if config.LogLevel == "" {
		config.LogLevel = *logLevel
	}

	_, exists := os.LookupEnv("STORE_INTERVAL")
	if !exists {
		config.StoreInterval = *storeInterval
	}

	if config.FileStoragePath == "" {
		config.FileStoragePath = *fileStoragePath
	}

	_, exists = os.LookupEnv("RESTORE")
	if !exists {
		config.Restore = *restore
	}

	return config, nil
}
