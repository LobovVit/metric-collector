package server

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/LobovVit/metric-collector/internal/server/config"
	"github.com/LobovVit/metric-collector/internal/server/domain/actions"
	"github.com/LobovVit/metric-collector/internal/server/server/middlewares"
	"github.com/LobovVit/metric-collector/pkg/logger"
	"github.com/LobovVit/metric-collector/pkg/postgresql"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	config  *config.Config
	storage actions.Repo
	dbCon   *sql.DB //*pgx.Conn
}

func New(config *config.Config) *Server {
	repo := actions.GetRepo(config.Restore, config.StoreInterval, config.FileStoragePath)
	return &Server{config: config, storage: repo}
}

func (a *Server) Run(ctx context.Context) error {

	dbCon, err := postgresql.NweConn(ctx, a.config.DSN)
	if err != nil {
		logger.Log.Error("Get db connection failed", zap.Error(err))
	}
	a.dbCon = dbCon

	mux := chi.NewRouter()
	mux.Use(middlewares.WithLogging)
	mux.Use(middlewares.WithCompress)
	mux.Get("/", a.allMetricsHandler)
	mux.Get("/ping/", a.dbPingHandler)
	mux.Post("/value/", a.singleMetricJSONHandler)
	mux.Get("/value/{type}/{name}", a.singleMetricHandler)
	mux.Post("/update/", a.updateJSONHandler)
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
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		logger.Log.Info("Shutdown", zap.Error(err))
		a.RouterShutdown()
	}
	return nil
}

func (a *Server) RouterShutdown() {
	err := a.storage.SaveToFile()
	if err != nil {
		logger.Log.Error("Save to file failed", zap.Error(err))
	}
}
