package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type configJSON struct {
	Address       string `json:"address,omitempty"`
	Restore       bool   `json:"restore,omitempty"`
	StoreInterval string `json:"store_interval,omitempty"`
	StoreFile     string `json:"store_file,omitempty"`
	DatabaseDsn   string `json:"database_dsn,omitempty"`
	CryptoKey     string `json:"crypto_key,omitempty"`
	TrustedSubnet string `json:"trusted_subnet,omitempty"`
}

func parseJSONConfig(config Config) (*Config, error) {
	cfgJSON := &configJSON{}
	if config.ConfigPath != "" {
		jsonFile, err := os.Open(config.ConfigPath)
		if err != nil {
			return &config, fmt.Errorf("JSON config file open: %w", err)
		}
		defer jsonFile.Close()
		byteValue, _ := io.ReadAll(jsonFile)
		err = json.Unmarshal(byteValue, cfgJSON)
		if err != nil {
			return &config, fmt.Errorf("JSON config unmarshal: %w", err)
		}
	}

	if config.Host == "" {
		config.Host = cfgJSON.Address
	}

	if config.CryptoKey == "" {
		config.CryptoKey = cfgJSON.CryptoKey
	}

	if config.StoreInterval == 0 {
		dur, err := time.ParseDuration(cfgJSON.StoreInterval)
		if err == nil {
			config.StoreInterval = int(dur.Seconds())
		}
	}
	if config.StoreInterval == 0 {
		config.StoreInterval = 30 //default value
	}

	if !config.Restore {
		config.Restore = cfgJSON.Restore
	}

	if config.FileStoragePath == "" {
		config.FileStoragePath = cfgJSON.StoreFile
	}

	if config.DSN == "" {
		config.DSN = cfgJSON.DatabaseDsn
	}

	if config.TrustedSubnet == "" {
		config.TrustedSubnet = cfgJSON.TrustedSubnet
	}

	return &config, nil
}
