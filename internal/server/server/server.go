package server

import (
	"context"
	"net/http"

	"github.com/LobovVit/metric-collector/internal/server/config"
	"github.com/LobovVit/metric-collector/internal/server/domain/actions"
	"github.com/LobovVit/metric-collector/internal/server/server/middleware"
	"github.com/LobovVit/metric-collector/pkg/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	config  *config.Config
	storage actions.Repo
}

func New(ctx context.Context, config *config.Config) (*Server, error) {
	repo, err := actions.GetRepo(ctx, config)
	if err != nil {
		return nil, err
	}
	return &Server{config: config, storage: repo}, nil
}

func (a *Server) Run(ctx context.Context) error {

	mux := chi.NewRouter()
	mux.Use(middleware.WithLogging)

	mux.With(middleware.WithCompress).
		Get("/", a.allMetricsHandler)
	mux.Get("/ping", a.dbPingHandler)
	mux.Post("/value/", a.singleMetricJSONHandler)
	mux.Get("/value/{type}/{name}", a.singleMetricHandler)
	mux.With(middleware.WithSignature(a.config.SigningKey)).
		With(middleware.WithCompress).
		Post("/update/", a.updateJSONHandler)
	mux.With(middleware.WithSignature(a.config.SigningKey)).
		With(middleware.WithCompress).
		Post("/updates/", a.updateBatchJSONHandler)
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

func (a *Server) RouterShutdown(ctx context.Context) {
	err := a.storage.SaveToFile(ctx)
	if err != nil {
		logger.Log.Error("Save to file failed", zap.Error(err))
	}
}
