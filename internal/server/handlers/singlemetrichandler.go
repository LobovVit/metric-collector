package handlers

import (
	"github.com/LobovVit/metric-collector/internal/server/domain/actions"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func SingleMetricHandler(w http.ResponseWriter, r *http.Request) {
	tp := strings.ToLower(chi.URLParam(r, "type"))
	name := strings.ToLower(chi.URLParam(r, "name"))
	w.Header().Set("Content-Type", "text/plain")
	res, err := actions.GetSingleVal(tp, name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res))
	}
}
