// Package middleware - included middleware for http handlers
package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/LobovVit/metric-collector/pkg/signature"
)

type signWriter struct {
	w    http.ResponseWriter
	hash []byte
	key  string
}

func newSignWriter(w http.ResponseWriter, key string) *signWriter {
	return &signWriter{
		w:   w,
		key: key,
	}
}

func (s *signWriter) WriteHeader(statusCode int) {
	s.w.Header().Set("HashSHA256", fmt.Sprintf("%x", s.hash))
	s.w.WriteHeader(statusCode)
}

func (s *signWriter) Header() http.Header {
	return s.w.Header()
}

func (s *signWriter) Write(p []byte) (int, error) {
	var err error
	s.hash, err = signature.CreateSignature(p, s.key)
	if err != nil {
		return 0, err
	}
	return s.w.Write(p)
}

// WithSignature - middleware for checking signature requests
func WithSignature(key string) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sw := w
			headerHash := r.Header.Get("HashSHA256")
			if key != "" && headerHash != "" {
				sw = newSignWriter(w, key)
				body, err := io.ReadAll(r.Body)
				defer r.Body.Close()
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				r.Body = io.NopCloser(bytes.NewBuffer(body))
				err = signature.CheckSignature(body, headerHash, key)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
			next.ServeHTTP(sw, r)
		}
		return http.HandlerFunc(fn)
	}
}
