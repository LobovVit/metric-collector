// Package config - included struct and init function fow work with app configuration
package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
)

// Config determines the basic parameters of the agent's operation
type Config struct {
	Host            string `env:"ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" `
	Restore         bool   `env:"RESTORE"`
	DSN             string `env:"DATABASE_DSN" `
	SigningKey      string `env:"KEY"`
	CryptoKey       string `env:"CRYPTO_KEY"`
	ConfigPath      string `env:"CONFIG"`
}

// GetConfig - method creates a new configuration and sets values from environment variables and command line flags
func GetConfig() (*Config, error) {
	config := &Config{}
	err := env.Parse(config)
	if err != nil {
		return nil, fmt.Errorf("env parse: %w", err)
	}

	host := flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	logLevel := flag.String("l", "info", "log level")
	storeInterval := flag.Int("i", 0, "интервал сохранения на диск")
	fileStoragePath := flag.String("f", "/tmp/metrics-db.json", "файл для сохранения на диск")
	restore := flag.Bool("r", false, "загружать при старте данные из файла")
	dsn := flag.String("d", "", "строка подключения к БД") //postgresql://postgres:password@10.66.66.3:5432/postgres?sslmode=disable
	signingKey := flag.String("k", "", "ключ")
	cryptoKey := flag.String("crypto-key", "", "путь до файла с закрытым ключом") //private.pem
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
	if config.DSN == "" {
		config.DSN = *dsn
	}
	if config.SigningKey == "" {
		config.SigningKey = *signingKey
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
