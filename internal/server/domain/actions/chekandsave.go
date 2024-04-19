package actions

import (
	"context"
	"fmt"
	"strconv"

	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
	"github.com/LobovVit/metric-collector/pkg/logger"
	"github.com/LobovVit/metric-collector/pkg/retry"
	"go.uber.org/zap"
)

type badRequestErr struct {
	tp    string
	value string
}

func (e badRequestErr) Error() string {
	return fmt.Sprintf("bad request metric type:\"%v\" with value:\"%v\"", e.tp, e.value)
}

func (r *Repo) CheckAndSaveText(ctx context.Context, tp string, name string, value string) error {
	var ret error
	try := retry.New(3)
	switch tp {
	case "gauge":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return badRequestErr{tp, value}
		}
		for {
			ret = r.storage.SetGauge(ctx, name, v)
			if ret == nil || !r.storage.IsRetryable(ret) || !try.Run() {
				break
			}
		}
	case "counter":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return badRequestErr{tp, value}
		}
		for {
			ret = r.storage.SetCounter(ctx, name, v)
			if ret == nil || !r.storage.IsRetryable(ret) || !try.Run() {
				break
			}
		}
	default:
		return badRequestErr{tp, value}
	}
	if r.needImmediatelySave {
		err := r.SaveToFile(ctx)
		if err != nil {
			logger.Log.Error("Immediately save failed", zap.Error(err))
		}
	}
	return ret
}

func (r *Repo) CheckAndSaveStruct(ctx context.Context, metrics metrics.Metrics) (metrics.Metrics, error) {
	var ret error
	try := retry.New(3)
	switch metrics.MType {
	case "gauge":
		for {
			ret = r.storage.SetGauge(ctx, metrics.ID, *metrics.Value)
			if ret == nil || !r.storage.IsRetryable(ret) || !try.Run() {
				break
			}
		}
	case "counter":
		for {
			ret = r.storage.SetCounter(ctx, metrics.ID, *metrics.Delta)
			if ret == nil || !r.storage.IsRetryable(ret) || !try.Run() {
				break
			}
		}
		tmp, _ := r.storage.GetSingle(ctx, metrics.MType, metrics.ID)
		*metrics.Delta, _ = strconv.ParseInt(tmp, 10, 64)
	default:
		return metrics, badRequestErr{metrics.MType, metrics.ID}
	}
	if r.needImmediatelySave {
		err := r.SaveToFile(ctx)
		if err != nil {
			logger.Log.Error("Immediately save failed", zap.Error(err))
		}
	}
	return metrics, ret
}

func (r *Repo) CheckAndSaveBatch(ctx context.Context, metrics []metrics.Metrics) ([]metrics.Metrics, error) {
	err := retry.DoWithoutReturn(ctx, 3, r.storage.SetBatch, metrics, r.storage.IsRetryable)
	if err != nil {
		return metrics, err
	}
	if r.needImmediatelySave {
		err = r.SaveToFile(ctx)
		if err != nil {
			logger.Log.Error("Immediately save failed", zap.Error(err))
		}
	}
	return metrics, nil
}
