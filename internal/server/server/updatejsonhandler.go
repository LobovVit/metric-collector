package server

import (
	"bytes"
	"encoding/json"
	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
	"net/http"
)

func (a *App) updateJSONHandler(w http.ResponseWriter, r *http.Request) {

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

	metric, err = a.storage.CheckAndSaveStruct(metric)
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
