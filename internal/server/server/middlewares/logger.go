package middlewares

import (
	"net/http"
	"time"

	"github.com/LobovVit/metric-collector/pkg/logger"
	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
	}
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 200,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Log.Info("Handler log",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Durationp("duration", &duration),
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
		)
	}
	return http.HandlerFunc(logFn)
}
