package server

import (
	"net/http/httptest"

	"github.com/LobovVit/metric-collector/internal/server/config"
	"github.com/go-chi/chi/v5"
)

var TS = ServerRun()

func ServerRun() *httptest.Server {
	mux := chi.NewRouter()
	cfg := config.Config{Host: "localhost:8080", LogLevel: "info", StoreInterval: 100, FileStoragePath: "1.json", Restore: false}
	tst := New(&cfg)
	mux.Get("/", tst.allMetricsHandler)
	mux.Post("/value/", tst.singleMetricJSONHandler)
	mux.Get("/value/{type}/{name}", tst.singleMetricHandler)
	mux.Post("/update/", tst.updateJSONHandler)
	mux.Post("/update/{type}/{name}/{value}", tst.updateHandler)
	return httptest.NewServer(mux)
}
