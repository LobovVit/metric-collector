package server

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
)

func (a *Server) updateJSONHandler(w http.ResponseWriter, r *http.Request) {

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

	metric, err = a.storage.CheckAndSaveStruct(r.Context(), metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	resp, err := json.Marshal(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
