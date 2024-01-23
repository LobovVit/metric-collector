package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LobovVit/metric-collector/internal/agent/config"
	"github.com/LobovVit/metric-collector/internal/agent/logger"
	"github.com/LobovVit/metric-collector/internal/agent/metrics"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
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

func (a *Agent) RunAgent(ctx context.Context, logLevel string) error {
	if err := logger.Initialize(logLevel); err != nil {
		return fmt.Errorf("log initialize failed: %w", err)
	}
	m := metrics.GetMetricStruct()

	readTicker := time.NewTicker(time.Second * time.Duration(a.cfg.PollInterval))
	sendTicker := time.NewTicker(time.Second * time.Duration(a.cfg.ReportInterval))
	defer sendTicker.Stop()
	defer readTicker.Stop()

	for {
		select {
		case <-readTicker.C:
			m.GetMetrics()
			logger.Log.Info("Read")
		case <-sendTicker.C:
			tmp := m.CounterExecMemStats
			m.CounterExecMemStats = 0
			err := a.sendRequest(ctx, m)
			if err != nil {
				m.CounterExecMemStats = tmp
				logger.Log.Info("Err sendRequest", zap.Error(err))
			}
		case <-ctx.Done():
			logger.Log.Info("Shutdown")
			return nil
		}
	}
}

func (a *Agent) sendRequest(ctx context.Context, metrics *metrics.Metrics) error {
	switch a.cfg.ReportFormat {
	case "json":
		return a.sendRequestJSON(ctx, metrics)
	case "text":
		return a.sendRequestText(ctx, metrics)
	default:
		return fmt.Errorf("incorrect format")
	}
}

func (a *Agent) sendRequestText(ctx context.Context, metrics *metrics.Metrics) error {
	var val string
	for _, v := range metrics.Metrics {
		switch v.MType {
		case "gauge":
			val = strconv.FormatFloat(*v.Value, 'f', 10, 64)
		case "counter":
			val = strconv.FormatInt(*v.Delta, 10)
		}
		_, err := a.client.R().
			SetContext(ctx).
			SetHeader("Content-Type", "text/plain").
			Post(fmt.Sprintf("%v/%v/%v/%v", a.cfg.Host, v.MType, v.ID, val))
		if err != nil {
			return fmt.Errorf("send request failed: %w", err)
		}
	}
	return nil
}

func (a *Agent) sendRequestJSON(ctx context.Context, metrics *metrics.Metrics) error {
	for _, v := range metrics.Metrics {
		metric, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("marshal json failed: %w", err)
		}
		_, err = a.client.R().
			SetContext(ctx).
			SetHeader("Content-Type", "application/json").
			SetBody(metric).
			Post(a.cfg.Host)
		logger.Log.Info("Sent", zap.Binary("JSON", metric))
		if err != nil {
			return fmt.Errorf("send request json failed: %w", err)
		}
	}
	return nil
}
