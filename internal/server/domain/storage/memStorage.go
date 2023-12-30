package storage

import (
	"errors"
	"fmt"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func GetStorage() *MemStorage {
	return &MemStorage{make(map[string]float64), make(map[string]int64)} //Storage
}

func (ms *MemStorage) SetGauge(key string, val float64) error {
	ms.gauge[key] = val
	return nil
}

func (ms *MemStorage) SetCounter(key string, val int64) error {
	ms.counter[key] += val
	return nil
}

func (ms *MemStorage) GetAll() map[string]map[string]string {
	retCounter := make(map[string]string)
	for k, v := range ms.counter {
		retCounter[k] = fmt.Sprintf("%d", v)
	}
	retGauge := make(map[string]string)
	for k, v := range ms.gauge {
		retGauge[k] = fmt.Sprintf("%f", v)
	}
	ret := make(map[string]map[string]string)
	ret["counter"] = retCounter
	ret["gauge"] = retGauge
	return ret
}

func (ms *MemStorage) GetSingle(tp string, name string) (string, error) {
	switch tp {
	case "gauge":
		res, ok := ms.gauge[name]
		if ok {
			return fmt.Sprintf("%f", res), nil
		} else {
			return "", errors.New("NotFound")
		}

	case "counter":
		res, ok := ms.counter[name]
		if ok {
			return fmt.Sprintf("%d", res), nil
		} else {
			return "", errors.New("NotFound")
		}
	default:
		return "", errors.New("NotFound")
	}
}
