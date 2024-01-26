package actions

import (
	"github.com/LobovVit/metric-collector/internal/server/domain/storage"
	"github.com/LobovVit/metric-collector/internal/server/logger"
	"go.uber.org/zap"
	"time"
)

type Repo struct {
	storage         repository
	storeInterval   int
	fileStoragePath string
}

type repository interface {
	SetGauge(key string, val float64) error
	SetCounter(key string, val int64) error
	GetAll() map[string]map[string]string
	GetSingle(tp string, name string) (string, error)
	SaveToFile(filename string) error
	LoadFromFile(filename string) error
}

func GetRepo(filename string, needRestore bool, storeInterval int, fileStoragePath string) Repo {
	logger.Log.Info("GetRepo", zap.String("GetRepo", filename), zap.Bool("needRestore", needRestore))
	return Repo{storage: storage.NewStorage(filename, needRestore), storeInterval: storeInterval, fileStoragePath: fileStoragePath}
}

func (r *Repo) SaveToFile(filename string) error {
	return r.storage.SaveToFile(filename)
}

func (r *Repo) LoadFromFile(filename string) error {
	return r.storage.LoadFromFile(filename)
}

func (r *Repo) RunPeriodicSave(filename string) {
	if r.storeInterval == 0 {
		return
	}
	saveTicker := time.NewTicker(time.Second * time.Duration(r.storeInterval))
	go func() {
		for {
			<-saveTicker.C
			err := r.storage.SaveToFile(filename)
			if err != nil {
				logger.Log.Info("periodic save failed", zap.Error(err))
			}
		}
	}()
}
