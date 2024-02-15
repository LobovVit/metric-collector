package server

import (
	"net/http"
)

func (a *Server) dbPingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	err := a.storage.Ping(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ОК"))
}
