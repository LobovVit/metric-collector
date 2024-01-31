package server

import (
	"net/http/httptest"

	"github.com/LobovVit/metric-collector/internal/server/config"
	"github.com/go-chi/chi/v5"
)

var Ts *httptest.Server = Server_Run()

func Server_Run() *httptest.Server {
	mux := chi.NewRouter()
	cfg, _ := config.GetConfig()
	tst := New(cfg)
	mux.Post("/update/", tst.updateJSONHandler)
	return httptest.NewServer(mux)
}
