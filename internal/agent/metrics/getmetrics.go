package metrics

import (
	"math/rand"
	"runtime"
)

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type Metrics struct {
	Metrics             map[string]Metric
	CounterExecMemStats int64
}

func GetMetricStruct() *Metrics {
	return &Metrics{Metrics: make(map[string]Metric, 30), CounterExecMemStats: 0}
}

func (m *Metrics) GetMetrics() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	m.CounterExecMemStats += 1
	//runtime
	alloc := float64(mem.Alloc)
	m.Metrics["Alloc"] = Metric{ID: "Alloc", MType: "gauge", Value: &alloc}
	buckHashSys := float64(mem.BuckHashSys)
	m.Metrics["BuckHashSys"] = Metric{ID: "BuckHashSys", MType: "gauge", Value: &buckHashSys}
	frees := float64(mem.Frees)
	m.Metrics["Frees"] = Metric{ID: "Frees", MType: "gauge", Value: &frees}
	gCCPUFraction := float64(mem.GCCPUFraction)
	m.Metrics["GCCPUFraction"] = Metric{ID: "GCCPUFraction", MType: "gauge", Value: &gCCPUFraction}
	gCSys := float64(mem.GCSys)
	m.Metrics["GCSys"] = Metric{ID: "GCSys", MType: "gauge", Value: &gCSys}
	heapAlloc := float64(mem.HeapAlloc)
	m.Metrics["HeapAlloc"] = Metric{ID: "HeapAlloc", MType: "gauge", Value: &heapAlloc}
	heapIdle := float64(mem.HeapIdle)
	m.Metrics["HeapIdle"] = Metric{ID: "HeapIdle", MType: "gauge", Value: &heapIdle}
	heapInuse := float64(mem.HeapInuse)
	m.Metrics["HeapInuse"] = Metric{ID: "HeapInuse", MType: "gauge", Value: &heapInuse}
	heapObjects := float64(mem.HeapObjects)
	m.Metrics["HeapObjects"] = Metric{ID: "HeapObjects", MType: "gauge", Value: &heapObjects}
	heapReleased := float64(mem.HeapReleased)
	m.Metrics["HeapReleased"] = Metric{ID: "HeapReleased", MType: "gauge", Value: &heapReleased}
	heapSys := float64(mem.HeapSys)
	m.Metrics["HeapSys"] = Metric{ID: "HeapSys", MType: "gauge", Value: &heapSys}
	lastGC := float64(mem.LastGC)
	m.Metrics["LastGC"] = Metric{ID: "LastGC", MType: "gauge", Value: &lastGC}
	lookups := float64(mem.Lookups)
	m.Metrics["Lookups"] = Metric{ID: "Lookups", MType: "gauge", Value: &lookups}
	mCacheInuse := float64(mem.MCacheInuse)
	m.Metrics["MCacheInuse"] = Metric{ID: "MCacheInuse", MType: "gauge", Value: &mCacheInuse}
	mCacheSys := float64(mem.MCacheSys)
	m.Metrics["MCacheSys"] = Metric{ID: "MCacheSys", MType: "gauge", Value: &mCacheSys}
	mSpanInuse := float64(mem.MSpanInuse)
	m.Metrics["MSpanInuse"] = Metric{ID: "MSpanInuse", MType: "gauge", Value: &mSpanInuse}
	mSpanSys := float64(mem.MSpanSys)
	m.Metrics["MSpanSys"] = Metric{ID: "MSpanSys", MType: "gauge", Value: &mSpanSys}
	mallocs := float64(mem.Mallocs)
	m.Metrics["Mallocs"] = Metric{ID: "Mallocs", MType: "gauge", Value: &mallocs}
	nextGC := float64(mem.NextGC)
	m.Metrics["NextGC"] = Metric{ID: "NextGC", MType: "gauge", Value: &nextGC}
	numForcedGC := float64(mem.NumForcedGC)
	m.Metrics["NumForcedGC"] = Metric{ID: "NumForcedGC", MType: "gauge", Value: &numForcedGC}
	numGC := float64(mem.NumGC)
	m.Metrics["NumGC"] = Metric{ID: "NumGC", MType: "gauge", Value: &numGC}
	otherSys := float64(mem.OtherSys)
	m.Metrics["OtherSys"] = Metric{ID: "OtherSys", MType: "gauge", Value: &otherSys}
	pauseTotalNs := float64(mem.PauseTotalNs)
	m.Metrics["PauseTotalNs"] = Metric{ID: "PauseTotalNs", MType: "gauge", Value: &pauseTotalNs}
	stackInuse := float64(mem.StackInuse)
	m.Metrics["StackInuse"] = Metric{ID: "StackInuse", MType: "gauge", Value: &stackInuse}
	stackSys := float64(mem.StackSys)
	m.Metrics["StackSys"] = Metric{ID: "StackSys", MType: "gauge", Value: &stackSys}
	sys := float64(mem.Sys)
	m.Metrics["Sys"] = Metric{ID: "Sys", MType: "gauge", Value: &sys}
	totalAlloc := float64(mem.TotalAlloc)
	m.Metrics["TotalAlloc"] = Metric{ID: "TotalAlloc", MType: "gauge", Value: &totalAlloc}
	//RandomValue
	randomValue := rand.Float64()
	m.Metrics["RandomValue"] = Metric{ID: "RandomValue", MType: "gauge", Value: &randomValue}
	//counter
	delta := m.CounterExecMemStats
	m.Metrics["PollCount"] = Metric{ID: "PollCount", MType: "counter", Delta: &delta}
}
