package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/LobovVit/metric-collector/internal/agent/compress"
	"github.com/LobovVit/metric-collector/internal/agent/config"
	"github.com/LobovVit/metric-collector/internal/agent/metrics"
	"github.com/LobovVit/metric-collector/pkg/logger"
	"github.com/LobovVit/metric-collector/pkg/retry"
	"github.com/LobovVit/metric-collector/pkg/signature"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type Agent struct {
	cfg    *config.Config
	client *resty.Client
}

func New(config *config.Config) *Agent {
	agent := Agent{cfg: config, client: resty.New()}
	return &agent
}

func (a *Agent) Run(ctx context.Context) error {
	m := metrics.GetMetricStruct()
	//var wg sync.WaitGroup
	readTicker := time.NewTicker(time.Second * time.Duration(a.cfg.PollInterval))
	defer readTicker.Stop()

	//GetMetrics
	go func() {
		for {
			select {
			case <-readTicker.C:
				m.GetMetrics()
				logger.Log.Info("Read")
			case <-ctx.Done():
				logger.Log.Info("Shutdown")
			}
		}
	}()

	//SendMetrics
	sendTicker := time.NewTicker(time.Second * time.Duration(a.cfg.ReportInterval))
	defer sendTicker.Stop()
	sem := newSemaphore(a.cfg.RateLimit)
	for {
		select {
		case <-sendTicker.C:
			tmp := m.CounterExecMemStats.Load()
			err := a.sendRequestWithRetry(ctx, m, sem)
			m.CounterExecMemStats.Store(m.CounterExecMemStats.Load() - tmp)
			if err != nil {
				m.CounterExecMemStats.Store(tmp)
				logger.Log.Error("Send request failed", zap.Error(err))
			}
			logger.Log.Info("Sent")
		case <-ctx.Done():
			logger.Log.Info("Shutdown")
			return nil
		}
	}
}

func (a *Agent) sendRequestWithRetry(ctx context.Context, metrics *metrics.Metrics, sem *semaphore) error {
	var err error
	try := retry.New(3)
	for {
		err = a.sendRequest(ctx, metrics, sem)
		if err == nil || !try.Run() {
			break
		}
	}
	return err
}

func (a *Agent) sendRequest(ctx context.Context, metrics *metrics.Metrics, sem *semaphore) error {
	metrics.RwMutex.RLock()
	defer metrics.RwMutex.RUnlock()

	switch a.cfg.ReportFormat {
	case "json":
		return a.sendRequestJSON(ctx, metrics, sem)
	case "text":
		return a.sendRequestText(ctx, metrics, sem)
	case "batch":
		return a.sendRequestBatchJSON(ctx, metrics, sem)
	default:
		return fmt.Errorf("incorrect format")
	}
}

func (a *Agent) sendRequestText(ctx context.Context, met *metrics.Metrics, sem *semaphore) error {
	var ret error = nil
	var wg sync.WaitGroup
	for _, v := range met.Metrics {
		wg.Add(1)
		go func(met metrics.Metric, sem *semaphore) {
			err := func() error {
				var val string
				sem.acquire()
				switch met.MType {
				case "gauge":
					val = strconv.FormatFloat(*met.Value, 'f', 10, 64)
				case "counter":
					val = strconv.FormatInt(*met.Delta, 10)
				}
				_, err := a.client.R().
					SetContext(ctx).
					SetHeader("Content-Type", "text/plain").
					Post(fmt.Sprintf("%v%v/%v/%v", a.cfg.Host, met.MType, met.ID, val))
				if err != nil {

					sem.release()
					return fmt.Errorf("send: %w", err)
				}
				sem.release()
				return nil
			}()
			if err != nil {
				ret = fmt.Errorf("send request: %w", err)
			}
		}(v, sem)
		wg.Done()
	}
	wg.Wait()
	return ret
}

func (a *Agent) sendRequestJSON(ctx context.Context, met *metrics.Metrics, sem *semaphore) error {
	var ret error = nil
	var wg sync.WaitGroup
	for _, v := range met.Metrics {
		wg.Add(1)
		go func(val metrics.Metric, sem *semaphore) {
			err := func() error {
				sem.acquire()
				metric, err := json.Marshal(val)
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
					sem.release()
					return fmt.Errorf("send json: %w", err)
				}
				sem.release()
				return nil
			}()
			if err != nil {
				ret = fmt.Errorf("send request: %w", err)
			}
		}(v, sem)
		wg.Done()
	}
	wg.Wait()
	return ret
}

func (a *Agent) sendRequestBatchJSON(ctx context.Context, met *metrics.Metrics, sem *semaphore) error {
	var maxPart int = len(met.Metrics) / a.cfg.MaxCntInBatch
	var ret error = nil
	var wg sync.WaitGroup
	for part := 0; part <= maxPart; part++ {
		wg.Add(1)
		go func(p int, sem *semaphore) {
			err := func() error {
				sem.acquire()
				startPos := p * a.cfg.MaxCntInBatch
				endPos := p*a.cfg.MaxCntInBatch + a.cfg.MaxCntInBatch
				if endPos > len(met.Metrics) {

					endPos = len(met.Metrics)
				}
				data, err := json.Marshal(met.Metrics[startPos:endPos])
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
					sem.release()
					return fmt.Errorf("send request json: %w", err)
				}
				sem.release()
				return nil
			}()
			if err != nil {
				ret = fmt.Errorf("send request: %w", err)
			}
		}(part, sem)
	}

	return ret
}
