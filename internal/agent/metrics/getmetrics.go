package metrics

import (
	"math/rand"
	"runtime"
	"sync"

	gopsutil "github.com/shirou/gopsutil/v3/mem"
)

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type Metrics struct {
	Metrics             []Metric     `json:"Metrics,omitempty"`
	CounterExecMemStats int64        `json:"-"`
	RwMutex             sync.RWMutex `json:"-"`
}

func GetMetricStruct() *Metrics {
	return &Metrics{Metrics: make([]Metric, 32), CounterExecMemStats: 0}
}

func (m *Metrics) GetMetricsRuntime() {
	m.RwMutex.Lock()
	defer m.RwMutex.Unlock()

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	m.CounterExecMemStats += 1
	//counter
	delta := m.CounterExecMemStats
	m.Metrics[0] = Metric{ID: "PollCount", MType: "counter", Delta: &delta}
	//runtime
	alloc := float64(mem.Alloc)
	m.Metrics[1] = Metric{ID: "Alloc", MType: "gauge", Value: &alloc}
	buckHashSys := float64(mem.BuckHashSys)
	m.Metrics[2] = Metric{ID: "BuckHashSys", MType: "gauge", Value: &buckHashSys}
	frees := float64(mem.Frees)
	m.Metrics[3] = Metric{ID: "Frees", MType: "gauge", Value: &frees}
	gCCPUFraction := float64(mem.GCCPUFraction)
	m.Metrics[4] = Metric{ID: "GCCPUFraction", MType: "gauge", Value: &gCCPUFraction}
	gCSys := float64(mem.GCSys)
	m.Metrics[5] = Metric{ID: "GCSys", MType: "gauge", Value: &gCSys}
	heapAlloc := float64(mem.HeapAlloc)
	m.Metrics[6] = Metric{ID: "HeapAlloc", MType: "gauge", Value: &heapAlloc}
	heapIdle := float64(mem.HeapIdle)
	m.Metrics[7] = Metric{ID: "HeapIdle", MType: "gauge", Value: &heapIdle}
	heapInuse := float64(mem.HeapInuse)
	m.Metrics[8] = Metric{ID: "HeapInuse", MType: "gauge", Value: &heapInuse}
	heapObjects := float64(mem.HeapObjects)
	m.Metrics[9] = Metric{ID: "HeapObjects", MType: "gauge", Value: &heapObjects}
	heapReleased := float64(mem.HeapReleased)
	m.Metrics[10] = Metric{ID: "HeapReleased", MType: "gauge", Value: &heapReleased}
	heapSys := float64(mem.HeapSys)
	m.Metrics[11] = Metric{ID: "HeapSys", MType: "gauge", Value: &heapSys}
	lastGC := float64(mem.LastGC)
	m.Metrics[12] = Metric{ID: "LastGC", MType: "gauge", Value: &lastGC}
	lookups := float64(mem.Lookups)
	m.Metrics[13] = Metric{ID: "Lookups", MType: "gauge", Value: &lookups}
	mCacheInuse := float64(mem.MCacheInuse)
	m.Metrics[14] = Metric{ID: "MCacheInuse", MType: "gauge", Value: &mCacheInuse}
	mCacheSys := float64(mem.MCacheSys)
	m.Metrics[15] = Metric{ID: "MCacheSys", MType: "gauge", Value: &mCacheSys}
	mSpanInuse := float64(mem.MSpanInuse)
	m.Metrics[16] = Metric{ID: "MSpanInuse", MType: "gauge", Value: &mSpanInuse}
	mSpanSys := float64(mem.MSpanSys)
	m.Metrics[17] = Metric{ID: "MSpanSys", MType: "gauge", Value: &mSpanSys}
	mallocs := float64(mem.Mallocs)
	m.Metrics[18] = Metric{ID: "Mallocs", MType: "gauge", Value: &mallocs}
	nextGC := float64(mem.NextGC)
	m.Metrics[19] = Metric{ID: "NextGC", MType: "gauge", Value: &nextGC}
	numForcedGC := float64(mem.NumForcedGC)
	m.Metrics[20] = Metric{ID: "NumForcedGC", MType: "gauge", Value: &numForcedGC}
	numGC := float64(mem.NumGC)
	m.Metrics[21] = Metric{ID: "NumGC", MType: "gauge", Value: &numGC}
	otherSys := float64(mem.OtherSys)
	m.Metrics[22] = Metric{ID: "OtherSys", MType: "gauge", Value: &otherSys}
	pauseTotalNs := float64(mem.PauseTotalNs)
	m.Metrics[23] = Metric{ID: "PauseTotalNs", MType: "gauge", Value: &pauseTotalNs}
	stackInuse := float64(mem.StackInuse)
	m.Metrics[24] = Metric{ID: "StackInuse", MType: "gauge", Value: &stackInuse}
	stackSys := float64(mem.StackSys)
	m.Metrics[25] = Metric{ID: "StackSys", MType: "gauge", Value: &stackSys}
	sys := float64(mem.Sys)
	m.Metrics[26] = Metric{ID: "Sys", MType: "gauge", Value: &sys}
	totalAlloc := float64(mem.TotalAlloc)
	m.Metrics[27] = Metric{ID: "TotalAlloc", MType: "gauge", Value: &totalAlloc}
	//RandomValue
	randomValue := rand.Float64()
	m.Metrics[28] = Metric{ID: "RandomValue", MType: "gauge", Value: &randomValue}
}

func (m *Metrics) GetMetricsGops() {
	m.RwMutex.Lock()
	defer m.RwMutex.Unlock()

	gops, _ := gopsutil.VirtualMemory()
	m.CounterExecMemStats += 1
	//counter
	delta := m.CounterExecMemStats
	m.Metrics[0] = Metric{ID: "PollCount", MType: "counter", Delta: &delta}
	//gopsutil
	totalMemory := float64(gops.Total)
	m.Metrics[29] = Metric{ID: "TotalMemory", MType: "gauge", Value: &totalMemory}
	freeMemory := float64(gops.Free)
	m.Metrics[30] = Metric{ID: "FreeMemory", MType: "gauge", Value: &freeMemory}
	cpu := gops.UsedPercent
	m.Metrics[31] = Metric{ID: "CPUutilization1", MType: "gauge", Value: &cpu}
}
