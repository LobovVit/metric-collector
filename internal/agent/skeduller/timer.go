package skeduller

import (
	"github.com/LobovVit/metric-collector/internal/agent/metrics"
	"sync"
	"time"
)

func StartTimer(readTime int64, sendTime int64, endPoint string) {
	m := metrics.GetMetricStruct()
	wg := sync.WaitGroup{}
	rw := sync.Mutex{}
	wg.Add(1)
	go func() {
		for _ = range time.Tick(time.Second * time.Duration(readTime)) {
			rw.Lock()
			m.GetMetrics()
			//fmt.Printf("readTime:%v valAlloc:%v\n", now, m.Gauge["Alloc"])
			rw.Unlock()
		}
	}()
	wg.Add(1)
	go func() {
		for _ = range time.Tick(time.Second * time.Duration(sendTime)) {
			rw.Lock()
			//fmt.Printf("sendTime:%v valAlloc:%v PollCount:%v\n", now, m.Gauge["Alloc"], m.Counter["PollCount"])
			m.CounterExecMemStats = 0
			sendRequest(m, endPoint)
			rw.Unlock()
		}
	}()
	wg.Wait()
}
