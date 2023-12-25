package metrics

import (
	"math/rand"
	"runtime"
)

var Mem runtime.MemStats

type Metrics struct {
	Gauge               map[string]float64
	Counter             map[string]int64
	CounterExecMemStats int64
}

func GetMetricStruct() *Metrics {
	return &Metrics{make(map[string]float64), make(map[string]int64), 0} //Storage
}

func (m *Metrics) GetMetrics() {

	runtime.ReadMemStats(&Mem)
	m.CounterExecMemStats += 1

	//runtime
	m.Gauge["Alloc"] = float64(Mem.Alloc)
	m.Gauge["BuckHashSys"] = float64(Mem.BuckHashSys)
	m.Gauge["Frees"] = float64(Mem.Frees)
	m.Gauge["GCCPUFraction"] = float64(Mem.GCCPUFraction)
	m.Gauge["GCSys"] = float64(Mem.GCSys)
	m.Gauge["HeapAlloc"] = float64(Mem.HeapAlloc)
	m.Gauge["HeapIdle"] = float64(Mem.HeapIdle)
	m.Gauge["HeapInuse"] = float64(Mem.HeapInuse)
	m.Gauge["HeapObjects"] = float64(Mem.HeapObjects)
	m.Gauge["HeapReleased"] = float64(Mem.HeapReleased)
	m.Gauge["HeapSys"] = float64(Mem.HeapSys)
	m.Gauge["LastGC"] = float64(Mem.LastGC)
	m.Gauge["Lookups"] = float64(Mem.Lookups)
	m.Gauge["MCacheInuse"] = float64(Mem.MCacheInuse)
	m.Gauge["MCacheSys"] = float64(Mem.MCacheSys)
	m.Gauge["MSpanInuse"] = float64(Mem.MSpanInuse)
	m.Gauge["MSpanSys"] = float64(Mem.MSpanSys)
	m.Gauge["Mallocs"] = float64(Mem.Mallocs)
	m.Gauge["NextGC"] = float64(Mem.NextGC)
	m.Gauge["NumForcedGC"] = float64(Mem.NumForcedGC)
	m.Gauge["NumGC"] = float64(Mem.NumGC)
	m.Gauge["OtherSys"] = float64(Mem.OtherSys)
	m.Gauge["PauseTotalNs"] = float64(Mem.PauseTotalNs)
	m.Gauge["StackInuse"] = float64(Mem.StackInuse)
	m.Gauge["StackSys"] = float64(Mem.StackSys)
	m.Gauge["Sys"] = float64(Mem.Sys)
	m.Gauge["TotalAlloc"] = float64(Mem.TotalAlloc)
	//RandomValue
	m.Gauge["RandomValue"] = rand.Float64()
	//counter
	m.Counter["PollCount"] = m.CounterExecMemStats
}
