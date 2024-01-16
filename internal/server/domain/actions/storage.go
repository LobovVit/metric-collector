package actions

import "github.com/LobovVit/metric-collector/internal/server/domain/storage"

type Repo struct {
	storage repository
}

func GetRepo() Repo {
	return Repo{storage.NewStorage()}
}

type repository interface {
	SetGauge(key string, val float64) error
	SetCounter(key string, val int64) error
	GetAll() map[string]map[string]string
	GetSingle(tp string, name string) (string, error)
}
