package app

import (
	"context"
	"fmt"
	"github.com/LobovVit/metric-collector/internal/agent/config"
	"github.com/LobovVit/metric-collector/internal/agent/metrics"
	"log"
	"strconv"
	"time"
)

type Agent struct {
	cfg *config.Config
	ctx context.Context
}

func NewAgent(config *config.Config, context context.Context) *Agent {
	agent := Agent{config, context}
	return &agent
}

func (a *Agent) RunAgent() {
	m := metrics.GetMetricStruct()

	readTicker := time.NewTicker(time.Second * time.Duration(a.cfg.PollInterval))
	sendTicker := time.NewTicker(time.Second * time.Duration(a.cfg.ReportInterval))
	defer sendTicker.Stop()
	defer readTicker.Stop()

	for {
		select {
		case <-readTicker.C:
			m.GetMetrics()
			log.Printf("read")
		case <-sendTicker.C:
			tmp := m.CounterExecMemStats
			m.CounterExecMemStats = 0
			err := a.sendRequest(m)
			if err != nil {
				m.CounterExecMemStats = tmp
				log.Printf("err sendRequest %v", err)
			}
			log.Printf("send")
		case <-a.ctx.Done():
			log.Printf("shutdown")
			return
		}
	}
}

func (a *Agent) sendRequest(metrics *metrics.Metrics) error {
	var sender restyClient
	sender.new()
	for k, v := range metrics.Gauge {
		_, err := sender.client.R().
			SetContext(a.ctx).
			SetHeader("Content-Type", "text/plain").
			Post(fmt.Sprintf("%vgauge/%v/%v", a.cfg.Host, k, strconv.FormatFloat(v, 'f', 10, 64)))
		if err != nil {
			return err
		}
	}
	for k, v := range metrics.Counter {
		_, err := sender.client.R().
			SetContext(a.ctx).
			SetHeader("Content-Type", "text/plain").
			Post(fmt.Sprintf("%vcounter/%v/%v", a.cfg.Host, k, strconv.FormatInt(v, 10)))
		if err != nil {
			return err
		}
	}
	return nil
}
