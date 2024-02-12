package actions

import (
	"fmt"
	"strconv"

	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
	"github.com/LobovVit/metric-collector/internal/server/domain/retry"
	"github.com/LobovVit/metric-collector/pkg/logger"
	"go.uber.org/zap"
)

type badRequestErr struct {
	tp    string
	value string
}

func (e badRequestErr) Error() string {
	return fmt.Sprintf("bad request metric type:\"%v\" with value:\"%v\"", e.tp, e.value)
}

func (r *Repo) CheckAndSaveText(tp string, name string, value string) error {
	var ret error
	repeat := retry.New(3)
	switch tp {
	case "gauge":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return badRequestErr{tp, value}
		}
		ret = repeat.RunKVFloatParam(r.storage.SetGauge, name, v)
	case "counter":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return badRequestErr{tp, value}
		}
		ret = repeat.RunKVIntParam(r.storage.SetCounter, name, v)
	default:
		return badRequestErr{tp, value}
	}
	if r.needImmediatelySave {
		err := r.SaveToFile()
		if err != nil {
			logger.Log.Error("Immediately save failed", zap.Error(err))
		}
	}
	return ret
}

func (r *Repo) CheckAndSaveStruct(metrics metrics.Metrics) (metrics.Metrics, error) {
	var ret error
	repeat := retry.New(3)
	switch metrics.MType {
	case "gauge":
		ret = repeat.RunKVFloatParam(r.storage.SetGauge, metrics.ID, *metrics.Value)
	case "counter":
		ret = repeat.RunKVIntParam(r.storage.SetCounter, metrics.ID, *metrics.Delta)
		tmp, _ := r.storage.GetSingle(metrics.MType, metrics.ID)
		*metrics.Delta, _ = strconv.ParseInt(tmp, 10, 64)
	default:
		return metrics, badRequestErr{metrics.MType, metrics.ID}
	}
	if r.needImmediatelySave {
		err := r.SaveToFile()
		if err != nil {
			logger.Log.Error("Immediately save failed", zap.Error(err))
		}
	}
	return metrics, ret
}

func (r *Repo) CheckAndSaveBatch(metrics []metrics.Metrics) ([]metrics.Metrics, error) {
	repeat := retry.New(3)
	ret := repeat.RunMetricsParam(r.storage.SetBatch, metrics)
	if r.needImmediatelySave {
		err := r.SaveToFile()
		if err != nil {
			logger.Log.Error("Immediately save failed", zap.Error(err))
		}
	}
	return metrics, ret
}
