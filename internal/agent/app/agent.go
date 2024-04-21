package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/LobovVit/metric-collector/internal/agent/compress"
	"github.com/LobovVit/metric-collector/internal/agent/config"
	"github.com/LobovVit/metric-collector/internal/agent/metrics"
	"github.com/LobovVit/metric-collector/pkg/logger"
	"github.com/LobovVit/metric-collector/pkg/signature"
)

// Agent struct is used to create Agent with settings.
type Agent struct {
	cfg    *config.Config
	client *resty.Client
}

// New method creates a new Agent.
func New(config *config.Config) *Agent {
	agent := Agent{cfg: config, client: resty.New()}
	return &agent
}

func (a *Agent) Run(ctx context.Context) error {
	m := metrics.GetMetricStruct()
	logger.Log.Info("Run")
	var wg sync.WaitGroup
	readTicker := time.NewTicker(time.Second * time.Duration(a.cfg.PollInterval))
	defer readTicker.Stop()
	//GetMetrics
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-readTicker.C:
				m.GetMetrics()
				logger.Log.Info("Read")
			case <-ctx.Done():
				logger.Log.Info("Shutdown")
				return
			}
		}
	}()

	//SendMetrics
	sendTicker := time.NewTicker(time.Second * time.Duration(a.cfg.ReportInterval))
	defer sendTicker.Stop()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-sendTicker.C:
				tmp := m.CounterExecMemStats.Load()
				err := a.sendRequestWithRetry(ctx, m)
				m.CounterExecMemStats.Store(m.CounterExecMemStats.Load() - tmp)
				if err != nil {
					m.CounterExecMemStats.Store(tmp)
					logger.Log.Error("Send request failed", zap.Error(err))
				}
				logger.Log.Info("Sent")
			case <-ctx.Done():
				logger.Log.Info("Shutdown")
				return
			}
		}
	}()
	wg.Wait()
	return nil
}

func (a *Agent) sendRequestWithRetry(ctx context.Context, metrics *metrics.Metrics) error {
	a.client.SetRetryCount(3).SetRetryWaitTime(3 * time.Second)
	return a.sendRequest(ctx, metrics)
}

func (a *Agent) sendRequest(ctx context.Context, metrics *metrics.Metrics) error {
	metrics.RwMutex.RLock()
	defer metrics.RwMutex.RUnlock()

	switch a.cfg.ReportFormat {
	case "json":
		return a.sendRequestsJSON(ctx, metrics)
	case "text":
		return a.sendRequestsText(ctx, metrics)
	case "batch":
		return a.sendRequestsBatchJSON(ctx, metrics)
	default:
		return fmt.Errorf("incorrect format")
	}
}

func (a *Agent) sendRequestsText(ctx context.Context, met *metrics.Metrics) error {
	g := errgroup.Group{}
	g.SetLimit(a.cfg.RateLimit)
	for _, v := range met.Metrics {
		val := v
		g.Go(func() error {
			return a.sendSingleRequestText(ctx, val)
		})
	}
	if err := g.Wait(); err != nil {
		logger.Log.Info("send", zap.Error(err))
		return fmt.Errorf("send: %w", err)
	}
	return nil
}

func (a *Agent) sendSingleRequestText(ctx context.Context, singleMetric metrics.Metric) error {
	var val string
	switch singleMetric.MType {
	case "gauge":
		val = strconv.FormatFloat(*singleMetric.Value, 'f', 10, 64)
	case "counter":
		val = strconv.FormatInt(*singleMetric.Delta, 10)
	}
	_, err := a.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("%v%v/%v/%v", a.cfg.Host, singleMetric.MType, singleMetric.ID, val))
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}
	return nil
}

func (a *Agent) sendRequestsJSON(ctx context.Context, met *metrics.Metrics) error {
	g := errgroup.Group{}
	g.SetLimit(a.cfg.RateLimit)
	for _, v := range met.Metrics {
		val := v
		g.Go(func() error {
			return a.sendSingleRequestJSON(ctx, val)
		})
	}
	if err := g.Wait(); err != nil {
		logger.Log.Info("send", zap.Error(err))
		return fmt.Errorf("send: %w", err)
	}
	return nil
}

func (a *Agent) sendSingleRequestJSON(ctx context.Context, singleMetric metrics.Metric) error {
	metric, err := json.Marshal(singleMetric)
	if err != nil {

		return fmt.Errorf("marshal json: %w", err)
	}
	metric, err = compress.Compress(metric)
	if err != nil {
		return fmt.Errorf("compress json: %w", err)
	}
	if a.cfg.SigningKey != "" {
		sign, err := signature.CreateSignature(metric, a.cfg.SigningKey)
		if err != nil {
			return fmt.Errorf("create signature: %w", err)
		}
		a.client.SetHeader("HashSHA256", fmt.Sprintf("%x", sign))
	}
	_, err = a.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(metric).
		Post(a.cfg.Host)

	if err != nil {
		return fmt.Errorf("send json: %w", err)
	}
	return nil
}

func (a *Agent) sendRequestsBatchJSON(ctx context.Context, met *metrics.Metrics) error {
	var maxPart = len(met.Metrics) / a.cfg.MaxCntInBatch
	g := errgroup.Group{}
	g.SetLimit(a.cfg.RateLimit)
	for part := 0; part <= maxPart; part++ {
		startPos := part * a.cfg.MaxCntInBatch
		endPos := part*a.cfg.MaxCntInBatch + a.cfg.MaxCntInBatch
		if endPos > len(met.Metrics) {
			endPos = len(met.Metrics)
		}
		val := met.Metrics[startPos:endPos]
		g.Go(func() error {
			return a.sendSingleRequestBatchJSON(ctx, val)
		})
	}
	if err := g.Wait(); err != nil {
		logger.Log.Info("send", zap.Error(err))
		return fmt.Errorf("send: %w", err)
	}
	return nil
}

func (a *Agent) sendSingleRequestBatchJSON(ctx context.Context, singlePartMetric []metrics.Metric) error {
	data, err := json.Marshal(singlePartMetric)
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}
	data, err = compress.Compress(data)
	if err != nil {
		return fmt.Errorf("compress json: %w", err)
	}
	if a.cfg.SigningKey != "" {
		sign, err := signature.CreateSignature(data, a.cfg.SigningKey)
		if err != nil {
			return fmt.Errorf("create signature: %w", err)
		}
		a.client.SetHeader("HashSHA256", fmt.Sprintf("%x", sign))
	}
	_, err = a.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(data).
		Post(a.cfg.Host)
	if err != nil {
		return fmt.Errorf("send request json: %w", err)
	}
	return nil
}
