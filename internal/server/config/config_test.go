package config

import (
	"reflect"
	"strconv"
	"testing"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    *Config
		wantErr bool
	}{
		{name: "test get config", want: &Config{
			Host:            "localhost:8080",
			LogLevel:        "info",
			StoreInterval:   30,
			FileStoragePath: "/tmp/metrics-db.json",
			Restore:         true,
			DSN:             "",
			SigningKey:      "",
			CryptoKey:       "private.pem",
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ADDRESS", tt.want.Host)
			t.Setenv("LOG_LEVEL", tt.want.LogLevel)
			t.Setenv("STORE_INTERVAL", strconv.Itoa(tt.want.StoreInterval))
			t.Setenv("FILE_STORAGE_PATH", tt.want.FileStoragePath)
			t.Setenv("RESTORE", strconv.FormatBool(tt.want.Restore))
			t.Setenv("DATABASE_DSN", tt.want.SigningKey)
			t.Setenv("KEY", tt.want.SigningKey)
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
