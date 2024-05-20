package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type configJSON struct {
	Address        string `json:"address,omitempty"`
	ReportInterval string `json:"report_interval,omitempty"`
	PollInterval   string `json:"poll_interval,omitempty"`
	CryptoKey      string `json:"crypto_key,omitempty"`
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

	if config.Host == "" && cfgJSON.Address != "" {
		config.Host = cfgJSON.Address
		if config.ReportFormat != "batch" {
			config.Host = "http://" + config.Host + "/update/"
		}
		if config.ReportFormat == "batch" {
			config.Host = "http://" + config.Host + "/updates/"
		}
	}

	if config.CryptoKey == "" {
		config.CryptoKey = cfgJSON.CryptoKey
	}

	if config.ReportInterval == 0 {
		dur, err := time.ParseDuration(cfgJSON.ReportInterval)
		if err == nil {
			config.ReportInterval = int64(dur.Seconds())
		}
	}
	if config.ReportInterval == 0 {
		config.ReportInterval = 10 //default value
	}

	if config.PollInterval == 0 {
		dur, err := time.ParseDuration(cfgJSON.PollInterval)
		if err == nil {
			config.PollInterval = int64(dur.Seconds())
		}
	}
	if config.PollInterval == 0 {
		config.PollInterval = 2 //default value
	}
	return &config, nil
}
