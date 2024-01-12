package skeduller

import (
	"context"
	"github.com/LobovVit/metric-collector/internal/agent/metrics"
	"log"
	"time"
)

func StartTimer(ctx context.Context, readTime int64, sendTime int64, endPoint string) {
	m := metrics.GetMetricStruct()

	readTicker := time.NewTicker(time.Second * time.Duration(readTime))
	sendTicker := time.NewTicker(time.Second * time.Duration(sendTime))
	defer sendTicker.Stop()
	defer readTicker.Stop()

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
				return
			}
			log.Printf("send\n")
		case <-ctx.Done():
			log.Printf("shutdown\n")
			return
		}
	}
}
