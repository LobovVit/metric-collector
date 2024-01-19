package server

import (
	"fmt"
	"github.com/LobovVit/metric-collector/internal/server/domain/actions"
	"github.com/LobovVit/metric-collector/internal/server/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
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

func (a *App) RouterRun(logLevel string) error {
	if err := logger.Initialize(logLevel); err != nil {
		return fmt.Errorf("log initialize failed: %w", err)
	}

	mux := chi.NewRouter()
	mux.Use(logger.WithLogging)

	mux.Get("/", a.allMetricsHandler)
	mux.Get("/value/{type}/{name}", a.singleMetricHandler)
	mux.Post("/update/{type}/{name}/{value}", a.updateHandler)

	logger.Log.Info("Running server", zap.String("address", a.host))
	return http.ListenAndServe(a.host, mux)
}
