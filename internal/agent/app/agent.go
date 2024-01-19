package app

import (
	"context"
	"fmt"
	"github.com/LobovVit/metric-collector/internal/agent/config"
	"github.com/LobovVit/metric-collector/internal/agent/metrics"
	"github.com/go-resty/resty/v2"
	"log"
	"strconv"
	"time"
)

type Agent struct {
	cfg    *config.Config
	client *resty.Client
}

func NewAgent(config *config.Config) *Agent {
	agent := Agent{cfg: config, client: resty.New()}
	agent.client.R().SetHeader("Content-Type", "text/plain")
	return &agent
}

func (a *Agent) RunAgent(ctx context.Context) {
	m := metrics.GetMetricStruct()

	readTicker := time.NewTicker(time.Second * time.Duration(a.cfg.PollInterval))
	sendTicker := time.NewTicker(time.Second * time.Duration(a.cfg.ReportInterval))
	defer sendTicker.Stop()
	defer readTicker.Stop()

	for {
		select {
		case <-readTicker.C:
			m.GetMetrics()
			log.Printf("Read")
		case <-sendTicker.C:
			tmp := m.CounterExecMemStats
			m.CounterExecMemStats = 0
			err := a.sendRequest(ctx, m)
			if err != nil {
				m.CounterExecMemStats = tmp
				log.Printf("Err sendRequest %v", err)
			}
			log.Printf("Sent")
		case <-ctx.Done():
			log.Printf("Shutdown")
			return
		}
	}
}

func (a *Agent) sendRequest(ctx context.Context, metrics *metrics.Metrics) error {
	for k, v := range metrics.Gauge {
		_, err := a.client.R().
			SetContext(ctx).
			Post(fmt.Sprintf("%vgauge/%v/%v", a.cfg.Host, k, strconv.FormatFloat(v, 'f', 10, 64)))
		if err != nil {
			return fmt.Errorf("send request failed: %w", err)
		}
	}
	for k, v := range metrics.Counter {
		_, err := a.client.R().
			SetContext(ctx).
			Post(fmt.Sprintf("%vcounter/%v/%v", a.cfg.Host, k, strconv.FormatInt(v, 10)))
		if err != nil {
			return fmt.Errorf("send request failed: %w", err)
		}
	}
	return nil
}
