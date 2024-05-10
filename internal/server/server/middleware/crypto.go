package middleware

import (
	"bytes"
	"io"
	"net/http"

	cryptorsa "github.com/LobovVit/metric-collector/pkg/crypto"
)

func RsaMiddleware(filepath string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err = r.Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			priv, err := cryptorsa.LoadPrivateKey(filepath)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			bodyDecode, err := cryptorsa.DecryptOAEP(priv, body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(bodyDecode))
			h.ServeHTTP(w, r)
		})
	}
}
