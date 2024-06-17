package config

import (
	"reflect"
	"strconv"
	"testing"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		want    *Config
		wantErr bool
	}{
		{name: "test get config", cfg: &Config{
			Host:           "localhost:8080",
			HostGRPC:       "localhost:3200",
			Mode:           "http",
			LogLevel:       "info",
			ReportInterval: 1,
			PollInterval:   1,
			ReportFormat:   "batch",
			SigningKey:     "",
			RateLimit:      10,
			MaxCntInBatch:  5,
		}, want: &Config{
			Host:           "http://localhost:8080/updates/",
			HostGRPC:       "localhost:3200",
			Mode:           "http",
			LogLevel:       "info",
			ReportInterval: 1,
			PollInterval:   1,
			ReportFormat:   "batch",
			SigningKey:     "",
			RateLimit:      10,
			MaxCntInBatch:  5,
			CryptoKey:      "public.pem",
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ADDRESS", tt.cfg.Host)
			t.Setenv("REPORT_INTERVAL", strconv.FormatInt(tt.cfg.ReportInterval, 10))
			t.Setenv("POLL_INTERVAL", strconv.FormatInt(tt.cfg.ReportInterval, 10))
			t.Setenv("LOG_LEVEL", tt.cfg.LogLevel)
			t.Setenv("REPORT_FORMAT", tt.cfg.ReportFormat)
			t.Setenv("KEY", tt.cfg.SigningKey)
			t.Setenv("RATE_LIMIT", strconv.Itoa(tt.cfg.RateLimit))
			t.Setenv("BATCH_LIMIT", strconv.Itoa(tt.cfg.MaxCntInBatch))
			t.Setenv("CRYPTO_KEY", tt.want.CryptoKey)
			got, err := GetConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
