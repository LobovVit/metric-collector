package app

import (
	"context"
	"fmt"

	"github.com/LobovVit/metric-collector/internal/agent/metrics"
	"github.com/LobovVit/metric-collector/internal/proto"
	"github.com/LobovVit/metric-collector/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (a *Agent) sendRequestGrpc(ctx context.Context, metrics *metrics.Metrics) error {
	c := proto.NewUpdateServicesClient(a.clientGRPC)

	metrics.RwMutex.RLock()
	defer metrics.RwMutex.RUnlock()

	switch a.cfg.ReportFormat {
	case "json":
		return a.sendSingleParallel(ctx, metrics, c)
	case "text":
		return a.sendSingleParallel(ctx, metrics, c)
	case "batch":
		return a.sendBatchParallel(ctx, metrics, c)
	default:
		return fmt.Errorf("incorrect format")
	}
}

func (a *Agent) sendSingleParallel(ctx context.Context, metrics *metrics.Metrics, c proto.UpdateServicesClient) error {
	g := errgroup.Group{}
	g.SetLimit(a.cfg.RateLimit)
	for _, v := range metrics.Metrics {
		val := v
		g.Go(func() error {
			return a.sendGrpcSingle(ctx, val, c)
		})
	}
	if err := g.Wait(); err != nil {
		logger.Log.Info("send", zap.Error(err))
		return fmt.Errorf("send: %w", err)
	}
	return nil
}

func (a *Agent) sendGrpcSingle(ctx context.Context, metric metrics.Metric, c proto.UpdateServicesClient) error {
	val := proto.Metric{Id: metric.ID, Type: proto.MetricTypes(proto.MetricTypes_value[metric.MType])}
	switch metric.MType {
	case "counter":
		val.Delta = *metric.Delta
	case "gauge":
		val.Value = *metric.Value
	}
	resp, err := c.SingleMetric(ctx, &val)
	if err != nil {
		return fmt.Errorf("grpc single metric: %w", err)
	}
	if resp.Error != "" {
		return fmt.Errorf("grpc single metric: %v", resp.Error)
	}
	return nil
}

func (a *Agent) sendBatchParallel(ctx context.Context, met *metrics.Metrics, c proto.UpdateServicesClient) error {
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
			return a.sendGrpcBatch(ctx, val, c)
		})
	}
	if err := g.Wait(); err != nil {
		logger.Log.Info("send", zap.Error(err))
		return fmt.Errorf("send: %w", err)
	}
	return nil
}

func (a *Agent) sendGrpcBatch(ctx context.Context, singlePartMetric []metrics.Metric, c proto.UpdateServicesClient) error {
	val := proto.Metrics{}
	for _, m := range singlePartMetric {
		tmp := proto.Metric{Id: m.ID, Type: proto.MetricTypes(proto.MetricTypes_value[m.MType])}
		switch m.MType {
		case "counter":
			tmp.Delta = *m.Delta
		case "gauge":
			tmp.Value = *m.Value
		}
		val.Metrics = append(val.Metrics, &tmp)
	}
	resp, err := c.ButchMetrics(ctx, &val)
	if err != nil {
		return fmt.Errorf("grpc batch: %w", err)
	}
	if resp.Error != "" {
		return fmt.Errorf("grpc batch: %v", resp.Error)
	}
	return nil
}
