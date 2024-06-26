// Package actions - contains methods for working with abstract storage
package actions

import (
	"context"
	"strconv"

	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
)

// GetSingleValStruct - method returns single value from storage, using struct values
func (r *Repo) GetSingleValStruct(ctx context.Context, metrics metrics.Metrics) (metrics.Metrics, error) {
	switch metrics.MType {
	case "gauge":
		val, err := r.GetSingle(ctx, metrics.MType, metrics.ID)
		if err != nil {
			return metrics, err
		}
		valFl, _ := strconv.ParseFloat(val, 64)
		metrics.Value = &valFl
	case "counter":
		val, err := r.GetSingle(ctx, metrics.MType, metrics.ID)
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
