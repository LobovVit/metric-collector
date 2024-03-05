package server

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
)

func (a *Server) updateBatchJSONHandler(w http.ResponseWriter, r *http.Request) {
	var metricsBatch []metrics.Metrics
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
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
