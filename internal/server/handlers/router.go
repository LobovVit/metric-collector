package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RouterRun(host string) error {
	mux := chi.NewRouter()

	mux.Get("/", allMetricsHandler)
	mux.Get("/value/{type}/{name}", singleMetricHandler)
	mux.Post("/update/{type}/{name}/{value}", updateHandler)

	return http.ListenAndServe(host, mux)
}
