package actions

import (
	"fmt"
	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
	"github.com/LobovVit/metric-collector/internal/server/logger"
	"go.uber.org/zap"
	"strconv"
)

type badRequestErr struct {
	tp    string
	value string
}

func (e badRequestErr) Error() string {
	return fmt.Sprintf("bad request metric type:\"%v\" with value:\"%v\"", e.tp, e.value)
}

func (r Repo) CheckAndSaveText(tp string, name string, value string) error {
	switch tp {
	case "gauge":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return badRequestErr{tp, value}
		}
		r.storage.SetGauge(name, v)
	case "counter":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return badRequestErr{tp, value}
		}
		r.storage.SetCounter(name, v)
	default:
		return badRequestErr{tp, value}
	}
	if r.storeInterval == 0 {
		err := r.SaveToFile(r.fileStoragePath)
		if err != nil {
			logger.Log.Info("immediately save failed", zap.Error(err))
		}
	}
	return nil
}

func (r Repo) CheckAndSaveStruct(metrics metrics.Metrics) (metrics.Metrics, error) {
	switch metrics.MType {
	case "gauge":
		r.storage.SetGauge(metrics.ID, *metrics.Value)
	case "counter":
		r.storage.SetCounter(metrics.ID, *metrics.Delta)
		tmp, _ := r.storage.GetSingle(metrics.MType, metrics.ID)
		*metrics.Delta, _ = strconv.ParseInt(tmp, 10, 64)
	default:
		return metrics, badRequestErr{metrics.MType, metrics.ID}
	}
	if r.storeInterval == 0 {
		err := r.SaveToFile(r.fileStoragePath)
		if err != nil {
			logger.Log.Info("immediately save failed", zap.Error(err))
		}
	}
	return metrics, nil
}
