package skeduller

import (
	"context"
	"github.com/LobovVit/metric-collector/internal/agent/metrics"
	"log"
	"sync"
	"time"
)

func StartTimer(ctx context.Context, readTime int64, sendTime int64, endPoint string) {
	m := metrics.GetMetricStruct()
	wg := sync.WaitGroup{}

	readTicker := time.NewTicker(time.Second * time.Duration(readTime))
	sendTicker := time.NewTicker(time.Second * time.Duration(sendTime))
	defer sendTicker.Stop()
	defer readTicker.Stop()

	wg.Add(1)
	for {
		select {
		case <-readTicker.C:
			m.GetMetrics()
			log.Printf("read\n")
		case <-sendTicker.C:
			m.CounterExecMemStats = 0
			err := sendRequest(ctx, m, endPoint)
			if err != nil {
				log.Printf("err sendRequest %v\n", err)
				wg.Done()
				return
			}
			log.Printf("send\n")
		case <-ctx.Done():
			log.Printf("shutdown\n")
			wg.Done()
			return
		}
	}
	wg.Wait()
}
