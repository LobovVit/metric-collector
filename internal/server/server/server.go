// Package server - included methods for running the http server, register handlers and middleware, and their implementation
package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/LobovVit/metric-collector/internal/server/config"
	"github.com/LobovVit/metric-collector/internal/server/domain/actions"
	"github.com/LobovVit/metric-collector/internal/server/server/middleware"
	"github.com/LobovVit/metric-collector/pkg/logger"
)

// Server - structure containing a server instance
type Server struct {
	config  *config.Config
	storage actions.Repo
}

// New - method to create server instance
func New(ctx context.Context, config *config.Config) (*Server, error) {
	repo, err := actions.GetRepo(ctx, config)
	if err != nil {
		return nil, err
	}
	return &Server{config: config, storage: repo}, nil
}

// Run - method to start server instance
func (a *Server) Run(ctx context.Context) error {

	mux := chi.NewRouter()

	mux.Use(middleware.WithLogging)
	mux.Use(middleware.WithSignature(a.config.SigningKey))
	mux.Use(middleware.WithCompress)
	if a.config.CryptoKey != "" {
		mux.Use(middleware.RsaMiddleware(a.config.CryptoKey))
	}

	mux.Get("/", a.allMetricsHandler)
	mux.Get("/ping", a.dbPingHandler)
	mux.Post("/value/", a.singleMetricJSONHandler)
	mux.Get("/value/{type}/{name}", a.singleMetricHandler)
	mux.Post("/update/", a.updateJSONHandler)
	mux.Post("/updates/", a.updateBatchJSONHandler)
	mux.Post("/update/{type}/{name}/{value}", a.updateHandler)

	logger.Log.Info("Starting server", zap.String("address", a.config.Host))

	httpServer := &http.Server{
		Addr:    a.config.Host,
		Handler: mux,
	}

	g := errgroup.Group{}
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-ctx.Done()
		return httpServer.Shutdown(ctx)
	})

	if err := g.Wait(); err != nil {
		logger.Log.Info("Shutdown", zap.Error(err))
		a.RouterShutdown(ctx)
	}
	return nil
}

// RouterShutdown - method that implements saving the server state when shutting down
func (a *Server) RouterShutdown(ctx context.Context) {
	err := a.storage.SaveToFile(ctx)
	if err != nil {
		logger.Log.Error("Save to file failed", zap.Error(err))
	}
}
