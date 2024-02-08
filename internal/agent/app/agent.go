package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/LobovVit/metric-collector/internal/agent/compress"
	"github.com/LobovVit/metric-collector/internal/agent/config"
	"github.com/LobovVit/metric-collector/internal/agent/metrics"
	"github.com/LobovVit/metric-collector/pkg/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type Agent struct {
	cfg    *config.Config
	client *resty.Client
}

func New(config *config.Config) *Agent {
	agent := Agent{cfg: config, client: resty.New()}
	agent.client.R().SetHeader("Content-Type", "text/plain")
	return &agent
}

func (a *Agent) Run(ctx context.Context) error {
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
				logger.Log.Error("Send request failed", zap.Error(err))
			}
			logger.Log.Info("Sent")
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
	case "batch":
		return a.sendRequestBatchJSON(ctx, metrics)
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
			Post(fmt.Sprintf("%v%v/%v/%v", a.cfg.Host, v.MType, v.ID, val))
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
		metric, err = compress.Compress(metric)
		if err != nil {
			return fmt.Errorf("compress json failed: %w", err)
		}
		_, err = a.client.R().
			SetContext(ctx).
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetBody(metric).
			Post(a.cfg.Host)

		if err != nil {
			return fmt.Errorf("send request json failed: %w", err)
		}
	}
	return nil
}

func (a *Agent) sendRequestBatchJSON(ctx context.Context, metrics *metrics.Metrics) error {
	data, err := json.Marshal(metrics.Metrics)
	if err != nil {
		return fmt.Errorf("marshal json failed: %w", err)
	}
	data, err = compress.Compress(data)
	if err != nil {
		return fmt.Errorf("compress json failed: %w", err)
	}
	_, err = a.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(data).
		Post(a.cfg.Host)

	if err != nil {
		return fmt.Errorf("send request json failed: %w", err)
	}
	return nil
}
