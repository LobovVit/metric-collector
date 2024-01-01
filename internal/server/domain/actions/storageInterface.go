package actions

import "github.com/LobovVit/metric-collector/internal/server/domain/storage"

var store = StorInterface(storage.GetStorage())

type StorInterface interface {
	SetGauge(key string, val float64) error
	SetCounter(key string, val int64) error
	GetAll() map[string]map[string]string
	GetSingle(tp string, name string) (string, error)
}
