package server

import (
	"bytes"
	"encoding/json"
	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
	"github.com/LobovVit/metric-collector/internal/server/logger"
	"go.uber.org/zap"
	"net/http"
)

func (a *App) singleMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	var metric metrics.Metrics
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info("req", zap.String("ID", metric.ID),
		zap.String("MType", metric.MType),
		zap.Int64("Delta", *metric.Delta),
		zap.Float64("MType", *metric.Value))

	res, err := a.storage.GetSingleValStruct(metric)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	logger.Log.Info("resp", zap.String("ID", res.ID),
		zap.String("MType", res.MType),
		zap.Int64("Delta", *res.Delta),
		zap.Float64("MType", *res.Value))
	resp, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
