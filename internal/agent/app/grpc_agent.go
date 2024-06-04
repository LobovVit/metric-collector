package app

import (
	"context"
	"fmt"

	"github.com/LobovVit/metric-collector/internal/agent/metrics"
	"github.com/LobovVit/metric-collector/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/LobovVit/metric-collector/proto"
)

func (a *Agent) sendRequestGrpc(ctx context.Context, metrics *metrics.Metrics) error {
	conn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err //TODO wrap
	}
	defer conn.Close()
	c := pb.NewUpdateServicesClient(conn)

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

func (a *Agent) sendSingleParallel(ctx context.Context, metrics *metrics.Metrics, c pb.UpdateServicesClient) error {
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

func (a *Agent) sendGrpcSingle(ctx context.Context, metric metrics.Metric, c pb.UpdateServicesClient) error {
	val := pb.Metric{Id: metric.ID, Type: pb.MetricTypes(pb.MetricTypes_value[metric.MType])}
	switch metric.MType {
	case "counter":
		val.Delta = *metric.Delta
	case "gauge":
		val.Value = *metric.Value
	}
	resp, err := c.SingleMetric(ctx, &val)
	if err != nil {
		return err
	}
	if resp.Error != "" {
		return err
	}
	return nil
}

func (a *Agent) sendBatchParallel(ctx context.Context, met *metrics.Metrics, c pb.UpdateServicesClient) error {
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

func (a *Agent) sendGrpcBatch(ctx context.Context, singlePartMetric []metrics.Metric, c pb.UpdateServicesClient) error {
	val := pb.Metrics{}
	for _, m := range singlePartMetric {
		tmp := pb.Metric{Id: m.ID, Type: pb.MetricTypes(pb.MetricTypes_value[m.MType])}
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
		return err
	}
	if resp.Error != "" {
		return err
	}
	return nil
}
