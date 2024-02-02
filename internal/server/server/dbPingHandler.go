package server

import (
	"context"
	"net/http"
)

func (a *Server) dbPingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	if a.dbCon == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err := a.dbCon.Ping(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ОК"))
}
