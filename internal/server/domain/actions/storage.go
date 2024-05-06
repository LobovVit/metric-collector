// Package actions - contains methods for working with abstract storage
package actions

import (
	"context"

	"github.com/LobovVit/metric-collector/internal/server/config"
	"github.com/LobovVit/metric-collector/internal/server/domain/dbstorage"
	"github.com/LobovVit/metric-collector/internal/server/domain/memstorage"
	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
	"github.com/LobovVit/metric-collector/pkg/retry"
)

// Repo - structure containing abstract storage
type Repo struct {
	storage
	needImmediatelySave bool
}

// storage - interface describes the behavior of the abstract storage
type storage interface {
	SetGauge(ctx context.Context, key string, val float64) error
	SetCounter(ctx context.Context, key string, val int64) error
	GetAll(ctx context.Context) (map[string]map[string]string, error)
	GetSingle(ctx context.Context, tp string, name string) (string, error)
	SaveToFile(ctx context.Context) error
	LoadFromFile(ctx context.Context) error
	Ping(ctx context.Context) error
	SetBatch(ctx context.Context, metrics []metrics.Metrics) error
	IsRetryable(err error) bool
}

// GetRepo - method returning a storage instance
func GetRepo(ctx context.Context, config *config.Config) (Repo, error) {
	if config.DSN == "" {
		nImmSave := false
		if config.StoreInterval == 0 {
			nImmSave = true
		}
		storage, err := memstorage.NewStorage(ctx, config.Restore, config.StoreInterval, config.FileStoragePath)
		if err != nil {
			return Repo{}, err
		}
		return Repo{storage: storage, needImmediatelySave: nImmSave}, nil
	}
	storage, err := dbstorage.NewStorage(ctx, config.DSN)
	if err != nil {
		return Repo{}, err
	}
	return Repo{storage: storage}, nil
}

// SaveToFile - method saves values from storage to file
func (r *Repo) SaveToFile(ctx context.Context) error {
	return retry.DoNoParams(ctx, 3, r.storage.SaveToFile, r.storage.IsRetryable)
}

// LoadFromFile - method loads values from file to storage
func (r *Repo) LoadFromFile(ctx context.Context) error {
	return retry.DoNoParams(ctx, 3, r.storage.LoadFromFile, r.storage.IsRetryable)
}

// Ping - method tests the connection to the database
func (r *Repo) Ping(ctx context.Context) error {
	return retry.DoNoParams(ctx, 3, r.storage.Ping, r.storage.IsRetryable)
}

// SetBatch - method writes array values to storage
func (r *Repo) SetBatch(ctx context.Context, metrics []metrics.Metrics) error {
	return retry.Do(ctx, 3, r.storage.SetBatch, metrics, r.storage.IsRetryable)
}
