// Package server - included methods for running the http server, register handlers and middleware, and their implementation
package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	cryptorsa "github.com/LobovVit/metric-collector/pkg/crypto"
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
	wg      sync.WaitGroup
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
		priv, err := cryptorsa.LoadPrivateKey(a.config.CryptoKey)
		if err != nil {
			mux.Use(middleware.RsaBad(err))
		}
		mux.Use(middleware.Rsa(priv))
	}

	if a.config.TrustedSubnet != "" {
		_, inet, err := net.ParseCIDR(a.config.TrustedSubnet)
		if err != nil {
			logger.Log.Error("cidr parse:", zap.Error(err))
		}
		if err == nil {
			mux.Use(middleware.WithCheckSubnet(inet))
		}
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

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})

	a.wg.Add(1)
	go (func() {
		<-gCtx.Done()
		a.Shutdown(httpServer)
	})()

	if err := g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) { //
		logger.Log.Error("http server:", zap.Error(err))
	}
	a.wg.Wait()
	return nil
}

// Shutdown - method that implements saving the server state when shutting down
func (a *Server) Shutdown(srv *http.Server) {
	defer a.wg.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := srv.Shutdown(shutdownCtx)
	if err != nil {
		logger.Log.Error("http server shutdown:", zap.Error(err))
	}
	if err == nil {
		logger.Log.Info("http server shutdown ok")
	}
	err = a.storage.SaveToFile(shutdownCtx)
	if err != nil {
		logger.Log.Error("Save to file:", zap.Error(err))
	}
	if err == nil {
		logger.Log.Info("Save to file ok")
	}
}
