package actions

import (
	"context"

	"github.com/LobovVit/metric-collector/internal/server/config"
	"github.com/LobovVit/metric-collector/internal/server/domain/dbstorage"
	"github.com/LobovVit/metric-collector/internal/server/domain/memstorage"
	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
	"github.com/LobovVit/metric-collector/pkg/retry"
)

type Repo struct {
	storage             repository
	needImmediatelySave bool
}

type repository interface {
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

func (r *Repo) SaveToFile(ctx context.Context) error {
	var err error
	try := retry.New(3)
	for {
		err = r.storage.SaveToFile(ctx)
		if err == nil || try.Run() || !r.storage.IsRetryable(err) {
			break
		}
	}
	return err
}

func (r *Repo) LoadFromFile(ctx context.Context) error {
	var err error
	try := retry.New(3)
	for {
		err = r.storage.LoadFromFile(ctx)
		if err == nil || try.Run() || !r.storage.IsRetryable(err) {
			break
		}
	}
	return err
}

func (r *Repo) Ping(ctx context.Context) error {
	var err error
	try := retry.New(3)
	for {
		err = r.storage.Ping(ctx)
		if err == nil || try.Run() || !r.storage.IsRetryable(err) {
			break
		}
	}
	return err
}

func (r *Repo) SetBatch(ctx context.Context, metrics []metrics.Metrics) error {
	var err error
	try := retry.New(3)
	for {
		err = r.storage.SetBatch(ctx, metrics)
		if err == nil || try.Run() || !r.storage.IsRetryable(err) {
			break
		}
	}
	return err
}
