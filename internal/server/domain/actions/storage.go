package actions

import (
	"github.com/LobovVit/metric-collector/internal/server/domain/memstorage"
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
}

func GetRepo(needRestore bool, storeInterval int, fileStoragePath string) Repo {
	nImmSave := false
	if storeInterval == 0 {
		nImmSave = true
	}
	return Repo{storage: memstorage.NewStorage(needRestore, storeInterval, fileStoragePath), needImmediatelySave: nImmSave}
}

func (r *Repo) SaveToFile() error {
	return r.storage.SaveToFile()
}

func (r *Repo) LoadFromFile() error {
	return r.storage.LoadFromFile()
}
