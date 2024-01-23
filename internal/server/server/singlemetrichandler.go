package server

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func (a *App) singleMetricHandler(w http.ResponseWriter, r *http.Request) {
	tp := strings.ToLower(chi.URLParam(r, "type"))
	name := strings.ToLower(chi.URLParam(r, "name"))
	w.Header().Set("Content-Type", "text/plain")
	res, err := a.storage.GetSingleValText(tp, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res))
	}
}
