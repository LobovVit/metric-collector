package server

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
	"github.com/LobovVit/metric-collector/pkg/logger"
	"go.uber.org/zap"
)

func (a *Server) updateBatchJSONHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("!!!!!updateBatchJSONHandler!!!!")
	var metricsBatch []metrics.Metrics
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Log.Info("!!!!!updateBatchJSONHandler!!!!", zap.Int("len", len(buf.Bytes())))
	if err = json.Unmarshal(buf.Bytes(), &metricsBatch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metricsBatch, err = a.storage.CheckAndSaveBatch(r.Context(), metricsBatch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	resp, err := json.Marshal(metricsBatch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
