package server

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strings"
)

func (a *App) updateHandler(w http.ResponseWriter, r *http.Request) {
	tp := strings.ToLower(chi.URLParam(r, "type"))
	name := strings.ToLower(chi.URLParam(r, "name"))
	value := strings.ToLower(chi.URLParam(r, "value"))
	log.Println("updateHandler:", tp, "/", name, "/", value)
	w.Header().Set("Content-Type", "text/plain")
	err := a.storage.CheckAndSave(tp, name, value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
