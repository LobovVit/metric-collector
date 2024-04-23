package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// updateHandler - handler that handles passing metrics through request parameters
func (a *Server) updateHandler(w http.ResponseWriter, r *http.Request) {
	tp := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")
	w.Header().Set("Content-Type", "text/plain")
	err := a.storage.CheckAndSaveText(r.Context(), tp, name, value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
