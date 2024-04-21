// Package logger - included functions for init logger
package logger

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log - singleton variable for logger
var (
	Log  = zap.NewNop()
	once sync.Once
)

// Initialize - function for initialize variable for logger
func Initialize(level string) error {
	var err error
	once.Do(func() {
		lvl, err := zap.ParseAtomicLevel(level)
		if err != nil {
			err = fmt.Errorf("log parse level: %w", err)
		}

		cfg := zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

		cfg.Level = lvl

		zl, err := cfg.Build()
		if err != nil {
			err = fmt.Errorf("log build: %w", err)
		}

		Log = zl
	})
	return err
}
