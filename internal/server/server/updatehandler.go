package server

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (a *App) updateHandler(w http.ResponseWriter, r *http.Request) {
	tp := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")
	w.Header().Set("Content-Type", "text/plain")
	err := a.storage.CheckAndSaveText(tp, name, value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
