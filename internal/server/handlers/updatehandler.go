package handlers

import (
	"github.com/LobovVit/metric-collector/internal/server/domain/actions"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func updateHandler(w http.ResponseWriter, r *http.Request) {
	tp := strings.ToLower(chi.URLParam(r, "type"))
	name := strings.ToLower(chi.URLParam(r, "name"))
	value := strings.ToLower(chi.URLParam(r, "value"))
	w.Header().Set("Content-Type", "text/plain")
	err := actions.CheckAndSave(tp, name, value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
