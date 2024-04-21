package memstorage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
	"github.com/LobovVit/metric-collector/pkg/logger"
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

func NewStorage(ctx context.Context, needRestore bool, storeInterval int, fileStoragePath string) (*MemStorage, error) {
	s := &MemStorage{Gauge: make(map[string]float64), Counter: make(map[string]int64), storeInterval: storeInterval, fileStoragePath: fileStoragePath}
	if needRestore {
		err := s.LoadFromFile(ctx)
		if err != nil {
			logger.Log.Error("Load from file failed", zap.Error(err))
		}
	}
	s.StartPeriodicSave(ctx)
	return s, nil
}

func (ms *MemStorage) SetGauge(ctx context.Context, key string, val float64) error {
	ms.rwGaugeMutex.Lock()
	defer ms.rwGaugeMutex.Unlock()
	ms.Gauge[key] = val
	return nil
}

func (ms *MemStorage) SetCounter(ctx context.Context, key string, val int64) error {
	ms.rwCounterMutex.Lock()
	defer ms.rwCounterMutex.Unlock()
	ms.Counter[key] += val
	return nil
}

func (ms *MemStorage) GetAll(ctx context.Context) (map[string]map[string]string, error) {
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
	return ret, nil
}

func (ms *MemStorage) GetSingle(ctx context.Context, tp string, name string) (string, error) {
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

func (ms *MemStorage) SaveToFile(ctx context.Context) error {
	ms.rwCounterMutex.RLock()
	defer ms.rwCounterMutex.RUnlock()
	ms.rwGaugeMutex.RLock()
	defer ms.rwGaugeMutex.RUnlock()

	tmpfile, err := os.Create(ms.fileStoragePath + "_tmp_")
	if err != nil {
		return fmt.Errorf("open tmp file: %w", err)
	}
	type tmpStorage struct {
		Gauge   map[string]float64 `json:"gauge"`
		Counter map[string]int64   `json:"counter"`
	}
	tmp := tmpStorage{Gauge: ms.Gauge, Counter: ms.Counter}
	data, err := json.MarshalIndent(tmp, "", "	")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	_, err = tmpfile.Write(data)
	if err != nil {
		return fmt.Errorf("write tmp: %w", err)
	}
	err = tmpfile.Close()
	if err != nil {
		return fmt.Errorf("close tmp: %w", err)
	}

	err = os.Rename(ms.fileStoragePath+"_tmp_", ms.fileStoragePath)
	if err != nil {
		return fmt.Errorf("rename file: %w", err)
	}
	return nil
}

func (ms *MemStorage) LoadFromFile(ctx context.Context) error {
	ms.rwCounterMutex.RLock()
	defer ms.rwCounterMutex.RUnlock()
	ms.rwGaugeMutex.RLock()
	defer ms.rwGaugeMutex.RUnlock()

	data, err := os.ReadFile(ms.fileStoragePath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	type tmpStorage struct {
		Gauge   map[string]float64 `json:"gauge"`
		Counter map[string]int64   `json:"counter"`
	}
	tmp := tmpStorage{}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	if len(tmp.Gauge) > 0 {
		ms.Gauge = tmp.Gauge
	}
	if len(tmp.Counter) > 0 {
		ms.Counter = tmp.Counter
	}
	return nil
}

func (ms *MemStorage) StartPeriodicSave(ctx context.Context) {
	if ms.storeInterval == 0 {
		return
	}
	saveTicker := time.NewTicker(time.Second * time.Duration(ms.storeInterval))
	go func() {
		for {
			<-saveTicker.C
			err := ms.SaveToFile(ctx)
			if err != nil {
				logger.Log.Error("Periodic save failed", zap.Error(err))
			}
		}
	}()
}

func (ms *MemStorage) Ping(ctx context.Context) error {
	return fmt.Errorf("no db")
}

func (ms *MemStorage) SetBatch(ctx context.Context, metrics []metrics.Metrics) error {
	ms.rwCounterMutex.Lock()
	defer ms.rwCounterMutex.Unlock()
	ms.rwGaugeMutex.Lock()
	defer ms.rwGaugeMutex.Unlock()

	for _, v := range metrics {
		if v.MType == "gauge" {
			ms.Gauge[v.ID] += *v.Value
		}
		if v.MType == "counter" {
			ms.Counter[v.ID] += *v.Delta
		}
	}
	return nil
}

func (ms *MemStorage) IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	var osErr *os.SyscallError
	return errors.As(err, &osErr)
}
