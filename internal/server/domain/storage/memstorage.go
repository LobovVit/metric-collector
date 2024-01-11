package storage

import (
	"errors"
	"fmt"
	"sync"
)

type MemStorage struct {
	Gauge          map[string]float64
	Counter        map[string]int64
	rwGaugeMutex   sync.RWMutex
	rwCounterMutex sync.RWMutex
}

func notFoundErr(tp, name string) error {
	return errors.New(fmt.Sprintf("Not Found metric type:\"%v\" with name:\"%v\"", tp, name))
}

func NewStorage() *MemStorage {
	return &MemStorage{make(map[string]float64), make(map[string]int64), sync.RWMutex{}, sync.RWMutex{}} //Storage
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
		} else {
			return "", notFoundErr(tp, name)
		}

	case "counter":
		ms.rwCounterMutex.RLock()
		defer ms.rwCounterMutex.RUnlock()

		res, ok := ms.Counter[name]
		if ok {
			return fmt.Sprintf("%d", res), nil
		} else {
			return "", notFoundErr(tp, name)
		}
	default:
		return "", notFoundErr(tp, name)
	}
}