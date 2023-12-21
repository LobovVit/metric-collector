package domain

import "fmt"

type DataInterface interface {
	SetGauge(key string, val float64) error
	SetCounter(key string, val int64) error
}

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

var Storage *MemStorage

func GetStorage() *MemStorage {
	Storage = &MemStorage{make(map[string]float64), make(map[string]int64)}
	return Storage
}

func (ms *MemStorage) SetGauge(key string, val float64) error {
	ms.gauge[key] = val
	return nil
}

func (ms *MemStorage) SetCounter(key string, val int64) error {
	ms.counter[key] += val
	return nil
}

func (ms *MemStorage) GetAll() {
	for k, v := range ms.counter {
		fmt.Printf("counter: %v=%v\n", k, v)
	}
	for k, v := range ms.gauge {
		fmt.Printf("gauge: %v=%v\n", k, v)
	}
}
