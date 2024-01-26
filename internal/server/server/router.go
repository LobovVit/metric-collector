package server

import (
	"context"
	"github.com/LobovVit/metric-collector/internal/server/compress"
	"github.com/LobovVit/metric-collector/internal/server/config"
	"github.com/LobovVit/metric-collector/internal/server/domain/actions"
	"github.com/LobovVit/metric-collector/internal/server/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type App struct {
	config  *config.Config
	storage actions.Repo
}

func GetApp(config *config.Config) *App {
	repo := actions.GetRepo(config.FileStoragePath, config.Restore, config.StoreInterval, config.FileStoragePath)
	return &App{config: config, storage: repo}
}

func (a *App) RouterRun(ctx context.Context) error {

	a.storage.RunPeriodicSave(a.config.FileStoragePath)

	mux := chi.NewRouter()
	mux.Use(logger.WithLogging)
	mux.Use(compress.WithCompress)
	mux.Get("/", a.allMetricsHandler)
	mux.Post("/value/", a.singleMetricJSONHandler)
	mux.Get("/value/{type}/{name}", a.singleMetricHandler)
	mux.Post("/update/", a.updateJSONHandler)
	mux.Post("/update/{type}/{name}/{value}", a.updateHandler)

	logger.Log.Info("main: starting server", zap.String("address", a.config.Host))

	httpServer := &http.Server{
		Addr:    a.config.Host,
		Handler: mux,
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		logger.Log.Info("shutdown", zap.Error(err))
		a.RouterShutdown()
	}
	return nil
}

func (a *App) RouterShutdown() {
	err := a.storage.SaveToFile(a.config.FileStoragePath)
	if err != nil {
		logger.Log.Info("save to file failed", zap.Error(err))
	}
}
