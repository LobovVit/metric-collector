package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// singleMetricHandler - handler processing the receipt of one metric through request parameters
func (a *Server) singleMetricHandler(w http.ResponseWriter, r *http.Request) {
	tp := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	w.Header().Set("Content-Type", "text/plain")
	res, err := a.storage.GetSingle(r.Context(), tp, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res))
	}
}
