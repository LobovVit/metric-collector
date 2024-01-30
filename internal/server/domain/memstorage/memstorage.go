package memstorage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/LobovVit/metric-collector/pkg/logger"
	"go.uber.org/zap"
)

type notFoundMetricError struct {
	tp   string
	name string
}

func (e notFoundMetricError) Error() string {
	return fmt.Sprintf("not found metric type:\"%v\" with name:\"%v\"", e.tp, e.name)
}

type MemStorage struct {
	Gauge           map[string]float64
	Counter         map[string]int64
	rwGaugeMutex    sync.RWMutex
	rwCounterMutex  sync.RWMutex
	storeInterval   int
	fileStoragePath string
}

func NewStorage(needRestore bool, storeInterval int, fileStoragePath string) *MemStorage {
	s := &MemStorage{Gauge: make(map[string]float64), Counter: make(map[string]int64), storeInterval: storeInterval, fileStoragePath: fileStoragePath}
	if needRestore {
		err := s.LoadFromFile()
		if err != nil {
			logger.Log.Error("Load from file failed", zap.Error(err))
		}
	}
	s.StartPeriodicSave()
	return s
}

func (ms *MemStorage) SetGauge(key string, val float64) error {
	ms.rwGaugeMutex.Lock()
	defer ms.rwGaugeMutex.Unlock()
	ms.Gauge[key] = val
	return nil
}

func (ms *MemStorage) SetCounter(key string, val int64) error {
	ms.rwCounterMutex.Lock()
	defer ms.rwCounterMutex.Unlock()
	ms.Counter[key] += val
	return nil
}

func (ms *MemStorage) GetAll() map[string]map[string]string {
	ms.rwCounterMutex.RLock()
	defer ms.rwCounterMutex.RUnlock()
	ms.rwGaugeMutex.RLock()
	defer ms.rwGaugeMutex.RUnlock()

	retCounter := make(map[string]string, len(ms.Counter))
	for k, v := range ms.Counter {
		retCounter[k] = fmt.Sprintf("%d", v)
	}
	retGauge := make(map[string]string, len(ms.Gauge))
	for k, v := range ms.Gauge {
		retGauge[k] = fmt.Sprintf("%f", v)
	}
	ret := make(map[string]map[string]string, 2)
	ret["counter"] = retCounter
	ret["gauge"] = retGauge
	return ret
}

func (ms *MemStorage) GetSingle(tp string, name string) (string, error) {
	switch tp {
	case "gauge":
		ms.rwGaugeMutex.RLock()
		defer ms.rwGaugeMutex.RUnlock()

		res, ok := ms.Gauge[name]
		if ok {
			return fmt.Sprintf("%g", res), nil
		}
	case "counter":
		ms.rwCounterMutex.RLock()
		defer ms.rwCounterMutex.RUnlock()

		res, ok := ms.Counter[name]
		if ok {
			return fmt.Sprintf("%d", res), nil
		}
	}
	return "", notFoundMetricError{tp, name}
}

func (ms *MemStorage) SaveToFile() error {
	ms.rwCounterMutex.RLock()
	defer ms.rwCounterMutex.RUnlock()
	ms.rwGaugeMutex.RLock()
	defer ms.rwGaugeMutex.RUnlock()

	tmpfile, err := os.Create(ms.fileStoragePath + "_tmp_")
	if err != nil {
		return fmt.Errorf("open tmp file failed: %w", err)
	}
	type tmpStorage struct {
		Gauge   map[string]float64 `json:"gauge"`
		Counter map[string]int64   `json:"counter"`
	}
	tmp := tmpStorage{Gauge: ms.Gauge, Counter: ms.Counter}
	data, err := json.MarshalIndent(tmp, "", "	")
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}
	_, err = tmpfile.Write(data)
	if err != nil {
		return fmt.Errorf("write tmp failed: %w", err)
	}
	err = tmpfile.Close()
	if err != nil {
		return fmt.Errorf("close tmp failed: %w", err)
	}

	err = os.Rename(ms.fileStoragePath+"_tmp_", ms.fileStoragePath)
	if err != nil {
		return fmt.Errorf("rename file failed: %w", err)
	}
	return nil
}

func (ms *MemStorage) LoadFromFile() error {
	ms.rwCounterMutex.RLock()
	defer ms.rwCounterMutex.RUnlock()
	ms.rwGaugeMutex.RLock()
	defer ms.rwGaugeMutex.RUnlock()

	data, err := os.ReadFile(ms.fileStoragePath)
	if err != nil {
		return fmt.Errorf("read file failed: %w", err)
	}
	type tmpStorage struct {
		Gauge   map[string]float64 `json:"gauge"`
		Counter map[string]int64   `json:"counter"`
	}
	tmp := tmpStorage{}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return fmt.Errorf("unmarshal failed: %w", err)
	}
	ms.Gauge = tmp.Gauge
	ms.Counter = tmp.Counter
	return nil
}

func (ms *MemStorage) StartPeriodicSave() {
	if ms.storeInterval == 0 {
		return
	}
	saveTicker := time.NewTicker(time.Second * time.Duration(ms.storeInterval))
	go func() {
		for {
			<-saveTicker.C
			err := ms.SaveToFile()
			if err != nil {
				logger.Log.Error("Periodic save failed", zap.Error(err))
			}
		}
	}()
}
