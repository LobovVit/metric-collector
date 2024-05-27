package middleware

import (
	"bytes"
	"crypto/rsa"
	"io"
	"net/http"

	cryptorsa "github.com/LobovVit/metric-collector/pkg/crypto"
)

func Rsa(priv *rsa.PrivateKey) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			err = r.Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			bodyDecode, err := cryptorsa.DecryptOAEP(priv, body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(bodyDecode))
			h.ServeHTTP(w, r)
		})
	}
}

func RsaBad(err error) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		})
	}
}
