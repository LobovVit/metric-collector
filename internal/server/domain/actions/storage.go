package actions

import (
	"github.com/LobovVit/metric-collector/internal/server/config"
	"github.com/LobovVit/metric-collector/internal/server/domain/dbstorage"
	"github.com/LobovVit/metric-collector/internal/server/domain/memstorage"
	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
)

type Repo struct {
	storage             repository
	needImmediatelySave bool
}

type repository interface {
	SetGauge(key string, val float64) error
	SetCounter(key string, val int64) error
	GetAll() map[string]map[string]string
	GetSingle(tp string, name string) (string, error)
	SaveToFile() error
	LoadFromFile() error
	Ping() error
	SetBatch(metrics []metrics.Metrics) error
}

func GetRepo(config *config.Config) Repo {
	if config.DSN == "" {
		nImmSave := false
		if config.StoreInterval == 0 {
			nImmSave = true
		}
		return Repo{storage: memstorage.NewStorage(config.Restore, config.StoreInterval, config.FileStoragePath), needImmediatelySave: nImmSave}
	}
	return Repo{storage: dbstorage.NewStorage(config.DSN)}
}

func (r *Repo) SaveToFile() error {
	return r.storage.SaveToFile()
}

func (r *Repo) LoadFromFile() error {
	return r.storage.LoadFromFile()
}

func (r *Repo) Ping() error {
	return r.storage.Ping()
}

func (r *Repo) SetBatch(metrics []metrics.Metrics) error {
	return r.storage.SetBatch(metrics)
}
