package server

import (
	"net/http"

	"github.com/LobovVit/metric-collector/pkg/logger"
)

func (a *Server) dbPingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	if a.dbCon == nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Log.Info("500 dbCon == nil")
		return
	}
	err := a.dbCon.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		logger.Log.Info("500 dbCon.Ping err")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ОК"))
}
