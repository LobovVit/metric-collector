package server

import (
	"github.com/LobovVit/metric-collector/internal/server/domain/actions"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type App struct {
	host    string
	storage actions.Repo
}

func GetApp(host string) *App {
	repo := actions.GetRepo()
	return &App{host: host, storage: repo}
}

func (a *App) RouterRun() error {
	mux := chi.NewRouter()

	mux.Get("/", a.allMetricsHandler)
	mux.Get("/value/{type}/{name}", a.singleMetricHandler)
	mux.Post("/update/{type}/{name}/{value}", a.updateHandler)

	return http.ListenAndServe(a.host, mux)
}
