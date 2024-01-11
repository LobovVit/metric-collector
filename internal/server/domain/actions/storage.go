package actions

import "github.com/LobovVit/metric-collector/internal/server/domain/storage"

type repo struct {
	storage repository
}

func getRepo() repository {
	return storage.NewStorage()
}

var store = repo{getRepo()}

type repository interface {
	SetGauge(key string, val float64) error
	SetCounter(key string, val int64) error
	GetAll() map[string]map[string]string
	GetSingle(tp string, name string) (string, error)
}
