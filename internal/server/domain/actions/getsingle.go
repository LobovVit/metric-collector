package actions

import (
	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
	"strconv"
)

func (r Repo) GetSingleValText(tp string, name string) (string, error) {
	return r.storage.GetSingle(tp, name)
}

func (r Repo) GetSingleValStruct(metrics metrics.Metrics) (metrics.Metrics, error) {
	switch metrics.MType {
	case "gauge":
		val, err := r.storage.GetSingle(metrics.MType, metrics.ID)
		if err != nil {
			return metrics, err
		}
		valFl, _ := strconv.ParseFloat(val, 64)
		metrics.Value = &valFl
	case "counter":
		val, err := r.storage.GetSingle(metrics.MType, metrics.ID)
		if err != nil {
			return metrics, err
		}
		valInt, _ := strconv.ParseInt(val, 10, 64)
		metrics.Delta = &valInt
	default:
		return metrics, badRequestErr{metrics.MType, metrics.ID}
	}
	return metrics, nil
}
