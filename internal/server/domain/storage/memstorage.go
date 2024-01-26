package storage

import (
	"encoding/json"
	"fmt"
	"github.com/LobovVit/metric-collector/internal/server/logger"
	"go.uber.org/zap"
	"os"
	"sync"
)

type notFoundMetricError struct {
	tp   string
	name string
}

func (e notFoundMetricError) Error() string {
	return fmt.Sprintf("not found metric type:\"%v\" with name:\"%v\"", e.tp, e.name)
}

type MemStorage struct {
	Gauge          map[string]float64
	Counter        map[string]int64
	rwGaugeMutex   sync.RWMutex
	rwCounterMutex sync.RWMutex
}

func NewStorage(filename string, needRestore bool) *MemStorage {
	s := &MemStorage{Gauge: make(map[string]float64), Counter: make(map[string]int64)}
	logger.Log.Info("NewStorage", zap.String("filename", filename), zap.Bool("needRestore", needRestore))
	if needRestore {
		err := s.LoadFromFile(filename)
		if err != nil {
			logger.Log.Info("LoadFromFile err", zap.Error(err))
		}
	}
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

func (ms *MemStorage) SaveToFile(filename string) error {
	ms.rwCounterMutex.RLock()
	defer ms.rwCounterMutex.RUnlock()
	ms.rwGaugeMutex.RLock()
	defer ms.rwGaugeMutex.RUnlock()

	tfile, err := os.Create(filename + "_tmp_")
	if err != nil {
		return fmt.Errorf("open tmp file failed: %w", err)
	}
	data, err := json.MarshalIndent(ms, "", "	")
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}
	_, err = tfile.Write(data)
	if err != nil {
		return fmt.Errorf("write tmp failed: %w", err)
	}
	tfile.Close()

	err = os.Rename(filename+"_tmp_", filename)
	if err != nil {
		return fmt.Errorf("rename file failed: %w", err)
	}
	return nil
}

func (ms *MemStorage) LoadFromFile(filename string) error {
	ms.rwCounterMutex.RLock()
	defer ms.rwCounterMutex.RUnlock()
	ms.rwGaugeMutex.RLock()
	defer ms.rwGaugeMutex.RUnlock()

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read file failed: %w", err)
	}
	//logger.Log.Info(filename, zap.String("data", string(data)))
	err = json.Unmarshal(data, ms)
	if err != nil {
		return fmt.Errorf("unmarshal failed: %w", err)
	}
	return nil
}
