package metrics

import (
	"math/rand"
	"runtime"
)

type Metrics struct {
	Gauge               map[string]float64
	Counter             map[string]int64
	CounterExecMemStats int64
}

func GetMetricStruct() *Metrics {
	return &Metrics{Gauge: make(map[string]float64), Counter: make(map[string]int64), CounterExecMemStats: 0} //Storage
}

func (m *Metrics) GetMetrics() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	m.CounterExecMemStats += 1
	//runtime
	m.Gauge["Alloc"] = float64(mem.Alloc)
	m.Gauge["BuckHashSys"] = float64(mem.BuckHashSys)
	m.Gauge["Frees"] = float64(mem.Frees)
	m.Gauge["GCCPUFraction"] = float64(mem.GCCPUFraction)
	m.Gauge["GCSys"] = float64(mem.GCSys)
	m.Gauge["HeapAlloc"] = float64(mem.HeapAlloc)
	m.Gauge["HeapIdle"] = float64(mem.HeapIdle)
	m.Gauge["HeapInuse"] = float64(mem.HeapInuse)
	m.Gauge["HeapObjects"] = float64(mem.HeapObjects)
	m.Gauge["HeapReleased"] = float64(mem.HeapReleased)
	m.Gauge["HeapSys"] = float64(mem.HeapSys)
	m.Gauge["LastGC"] = float64(mem.LastGC)
	m.Gauge["Lookups"] = float64(mem.Lookups)
	m.Gauge["MCacheInuse"] = float64(mem.MCacheInuse)
	m.Gauge["MCacheSys"] = float64(mem.MCacheSys)
	m.Gauge["MSpanInuse"] = float64(mem.MSpanInuse)
	m.Gauge["MSpanSys"] = float64(mem.MSpanSys)
	m.Gauge["Mallocs"] = float64(mem.Mallocs)
	m.Gauge["NextGC"] = float64(mem.NextGC)
	m.Gauge["NumForcedGC"] = float64(mem.NumForcedGC)
	m.Gauge["NumGC"] = float64(mem.NumGC)
	m.Gauge["OtherSys"] = float64(mem.OtherSys)
	m.Gauge["PauseTotalNs"] = float64(mem.PauseTotalNs)
	m.Gauge["StackInuse"] = float64(mem.StackInuse)
	m.Gauge["StackSys"] = float64(mem.StackSys)
	m.Gauge["Sys"] = float64(mem.Sys)
	m.Gauge["TotalAlloc"] = float64(mem.TotalAlloc)
	//RandomValue
	m.Gauge["RandomValue"] = rand.Float64()
	//counter
	m.Counter["PollCount"] = m.CounterExecMemStats
}
