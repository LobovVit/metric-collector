package middleware

import (
	"fmt"
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

func Signature(key string) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sw := w
			if key != "" {
				sw = newSignWriter(w, key)
				ok, err := signature.CheckSignature([]byte(r.Header.Get("HashSHA256")), key)
				if err != nil || !ok {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
			next.ServeHTTP(sw, r)
		}
		return http.HandlerFunc(fn)
	}
}
