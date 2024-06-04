// Package server - included methods for running the http server, register handlers and middleware, and their implementation
package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
	cryptorsa "github.com/LobovVit/metric-collector/pkg/crypto"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/LobovVit/metric-collector/internal/server/config"
	"github.com/LobovVit/metric-collector/internal/server/domain/actions"
	"github.com/LobovVit/metric-collector/internal/server/server/middleware"
	"github.com/LobovVit/metric-collector/pkg/logger"
	pb "github.com/LobovVit/metric-collector/proto"
)

// Server - structure containing a server instance
type Server struct {
	config  *config.Config
	storage actions.Repo
	wg      sync.WaitGroup
	pb.UnimplementedUpdateServicesServer
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

	httpServer := &http.Server{
		Addr:    a.config.Host,
		Handler: mux,
	}

	grpcListener, err := net.Listen("tcp", a.config.HostGRPC)
	if err != nil {
		return fmt.Errorf("grpc listen: %w", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUpdateServicesServer(grpcServer, a)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		logger.Log.Info("Starting http server", zap.String("address", a.config.Host))
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		logger.Log.Info("Starting grpc server", zap.String("address", a.config.HostGRPC))
		return grpcServer.Serve(grpcListener)
	})

	a.wg.Add(1)
	go (func() {
		<-gCtx.Done()
		a.Shutdown(httpServer, grpcServer)
	})()

	if err := g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) { //
		logger.Log.Error("http server:", zap.Error(err))
	}
	a.wg.Wait()
	return nil
}

// Shutdown - method that implements saving the server state when shutting down
func (a *Server) Shutdown(srvHTTP *http.Server, srvGRPC *grpc.Server) {
	defer a.wg.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := srvHTTP.Shutdown(shutdownCtx)
	if err != nil {
		logger.Log.Error("http server shutdown:", zap.Error(err))
	}
	if err == nil {
		logger.Log.Info("http server shutdown ok")
	}

	srvGRPC.GracefulStop()
	logger.Log.Info("grpc server shutdown ok")

	err = a.storage.SaveToFile(shutdownCtx)
	if err != nil {
		logger.Log.Error("Save to file:", zap.Error(err))
	}
	if err == nil {
		logger.Log.Info("Save to file ok")
	}
}

func (a *Server) SingleMetric(ctx context.Context, in *pb.Metric) (*pb.Response, error) {
	var response pb.Response
	metric := metrics.Metrics{ID: in.Id, MType: in.Type.String(), Delta: &in.Delta, Value: &in.Value}
	_, err := a.storage.CheckAndSaveStruct(ctx, metric)
	if err != nil {
		response.Error = fmt.Sprintf("single metric check and save: %v", err)
	}
	return &response, nil
}

func (a *Server) ButchMetrics(ctx context.Context, in *pb.Metrics) (*pb.Response, error) {
	var response pb.Response
	var sliceMetric []metrics.Metrics
	for _, metric := range in.Metrics {
		sliceMetric = append(sliceMetric, metrics.Metrics{ID: metric.Id, MType: metric.Type.String(), Delta: &metric.Delta, Value: &metric.Value})
	}
	_, err := a.storage.CheckAndSaveBatch(ctx, sliceMetric)
	if err != nil {
		response.Error = fmt.Sprintf("batch metrics check and save: %v", err)
	}
	return &response, nil
}
