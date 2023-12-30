package main

import (
	"github.com/LobovVit/metric-collector/internal/server/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	//mux := http.NewServeMux()
	mux := chi.NewRouter()
	//mux.HandleFunc(`/update/`, handlers.UpdateHandler)
	mux.Get("/", handlers.AllMetricsHandler)
	mux.Get("/value/{type}/{name}", handlers.SingleMetricHandler)
	mux.Post("/update/{type}/{name}/{value}", handlers.UpdateHandler)

	return http.ListenAndServe(`:8080`, mux)
}
