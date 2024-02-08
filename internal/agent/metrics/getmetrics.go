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
	Metrics             []Metric `json:"Metrics,omitempty"`
	CounterExecMemStats int64    `json:"-"`
}

func GetMetricStruct() *Metrics {
	return &Metrics{Metrics: nil, CounterExecMemStats: 0}
}

func (m *Metrics) GetMetrics() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	m.CounterExecMemStats += 1
	//runtime
	m.Metrics = nil
	alloc := float64(mem.Alloc)
	m.Metrics = append(m.Metrics, Metric{ID: "Alloc", MType: "gauge", Value: &alloc})
	buckHashSys := float64(mem.BuckHashSys)
	m.Metrics = append(m.Metrics, Metric{ID: "BuckHashSys", MType: "gauge", Value: &buckHashSys})
	frees := float64(mem.Frees)
	m.Metrics = append(m.Metrics, Metric{ID: "Frees", MType: "gauge", Value: &frees})
	gCCPUFraction := float64(mem.GCCPUFraction)
	m.Metrics = append(m.Metrics, Metric{ID: "GCCPUFraction", MType: "gauge", Value: &gCCPUFraction})
	gCSys := float64(mem.GCSys)
	m.Metrics = append(m.Metrics, Metric{ID: "GCSys", MType: "gauge", Value: &gCSys})
	heapAlloc := float64(mem.HeapAlloc)
	m.Metrics = append(m.Metrics, Metric{ID: "HeapAlloc", MType: "gauge", Value: &heapAlloc})
	heapIdle := float64(mem.HeapIdle)
	m.Metrics = append(m.Metrics, Metric{ID: "HeapIdle", MType: "gauge", Value: &heapIdle})
	heapInuse := float64(mem.HeapInuse)
	m.Metrics = append(m.Metrics, Metric{ID: "HeapInuse", MType: "gauge", Value: &heapInuse})
	heapObjects := float64(mem.HeapObjects)
	m.Metrics = append(m.Metrics, Metric{ID: "HeapObjects", MType: "gauge", Value: &heapObjects})
	heapReleased := float64(mem.HeapReleased)
	m.Metrics = append(m.Metrics, Metric{ID: "HeapReleased", MType: "gauge", Value: &heapReleased})
	heapSys := float64(mem.HeapSys)
	m.Metrics = append(m.Metrics, Metric{ID: "HeapSys", MType: "gauge", Value: &heapSys})
	lastGC := float64(mem.LastGC)
	m.Metrics = append(m.Metrics, Metric{ID: "LastGC", MType: "gauge", Value: &lastGC})
	lookups := float64(mem.Lookups)
	m.Metrics = append(m.Metrics, Metric{ID: "Lookups", MType: "gauge", Value: &lookups})
	mCacheInuse := float64(mem.MCacheInuse)
	m.Metrics = append(m.Metrics, Metric{ID: "MCacheInuse", MType: "gauge", Value: &mCacheInuse})
	mCacheSys := float64(mem.MCacheSys)
	m.Metrics = append(m.Metrics, Metric{ID: "MCacheSys", MType: "gauge", Value: &mCacheSys})
	mSpanInuse := float64(mem.MSpanInuse)
	m.Metrics = append(m.Metrics, Metric{ID: "MSpanInuse", MType: "gauge", Value: &mSpanInuse})
	mSpanSys := float64(mem.MSpanSys)
	m.Metrics = append(m.Metrics, Metric{ID: "MSpanSys", MType: "gauge", Value: &mSpanSys})
	mallocs := float64(mem.Mallocs)
	m.Metrics = append(m.Metrics, Metric{ID: "Mallocs", MType: "gauge", Value: &mallocs})
	nextGC := float64(mem.NextGC)
	m.Metrics = append(m.Metrics, Metric{ID: "NextGC", MType: "gauge", Value: &nextGC})
	numForcedGC := float64(mem.NumForcedGC)
	m.Metrics = append(m.Metrics, Metric{ID: "NumForcedGC", MType: "gauge", Value: &numForcedGC})
	numGC := float64(mem.NumGC)
	m.Metrics = append(m.Metrics, Metric{ID: "NumGC", MType: "gauge", Value: &numGC})
	otherSys := float64(mem.OtherSys)
	m.Metrics = append(m.Metrics, Metric{ID: "OtherSys", MType: "gauge", Value: &otherSys})
	pauseTotalNs := float64(mem.PauseTotalNs)
	m.Metrics = append(m.Metrics, Metric{ID: "PauseTotalNs", MType: "gauge", Value: &pauseTotalNs})
	stackInuse := float64(mem.StackInuse)
	m.Metrics = append(m.Metrics, Metric{ID: "StackInuse", MType: "gauge", Value: &stackInuse})
	stackSys := float64(mem.StackSys)
	m.Metrics = append(m.Metrics, Metric{ID: "StackSys", MType: "gauge", Value: &stackSys})
	sys := float64(mem.Sys)
	m.Metrics = append(m.Metrics, Metric{ID: "Sys", MType: "gauge", Value: &sys})
	totalAlloc := float64(mem.TotalAlloc)
	m.Metrics = append(m.Metrics, Metric{ID: "TotalAlloc", MType: "gauge", Value: &totalAlloc})
	//RandomValue
	randomValue := rand.Float64()
	m.Metrics = append(m.Metrics, Metric{ID: "RandomValue", MType: "gauge", Value: &randomValue})
	//counter
	delta := m.CounterExecMemStats
	m.Metrics = append(m.Metrics, Metric{ID: "PollCount", MType: "counter", Delta: &delta})
}
