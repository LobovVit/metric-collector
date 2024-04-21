package actions

import (
	"context"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
	"github.com/LobovVit/metric-collector/pkg/logger"
	"github.com/LobovVit/metric-collector/pkg/retry"
)

type badRequestErr struct {
	tp    string
	value string
}

func (e badRequestErr) Error() string {
	return fmt.Sprintf("bad request metric type:\"%v\" with value:\"%v\"", e.tp, e.value)
}

func (r *Repo) CheckAndSaveText(ctx context.Context, tp string, name string, value string) error {
	switch tp {
	case "gauge":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return badRequestErr{tp, value}
		}
		err = retry.DoTwoParams(ctx, 3, r.storage.SetGauge, name, v, r.storage.IsRetryable)
		if err != nil {
			return err
		}
	case "counter":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return badRequestErr{tp, value}
		}
		err = retry.DoTwoParams(ctx, 3, r.storage.SetCounter, name, v, r.storage.IsRetryable)
		if err != nil {
			return err
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
	return nil
}

func (r *Repo) CheckAndSaveStruct(ctx context.Context, metrics metrics.Metrics) (metrics.Metrics, error) {
	switch metrics.MType {
	case "gauge":
		err := retry.DoTwoParams(ctx, 3, r.storage.SetGauge, metrics.ID, *metrics.Value, r.storage.IsRetryable)
		if err != nil {
			return metrics, err
		}
	case "counter":
		err := retry.DoTwoParams(ctx, 3, r.storage.SetCounter, metrics.ID, *metrics.Delta, r.storage.IsRetryable)
		if err != nil {
			return metrics, err
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
	return metrics, nil
}

func (r *Repo) CheckAndSaveBatch(ctx context.Context, metrics []metrics.Metrics) ([]metrics.Metrics, error) {
	err := retry.Do(ctx, 3, r.storage.SetBatch, metrics, r.storage.IsRetryable)
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
