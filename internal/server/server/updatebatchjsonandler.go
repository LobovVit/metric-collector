package server

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
	"github.com/LobovVit/metric-collector/pkg/logger"
	"go.uber.org/zap"
)

func (a *Server) updateBatchJSONandler(w http.ResponseWriter, r *http.Request) {

	var metricsBatch []metrics.Metrics //metrics.SlMetrics
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	logger.Log.Info("body", zap.String("q", buf.String()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &metricsBatch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metricsBatch, err = a.storage.CheckAndSaveBatch(metricsBatch)
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
