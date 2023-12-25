package handlers

import (
	"github.com/LobovVit/metric-collector/internal/server/domain/actions"
	"net/http"
)

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		//только POST
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(actions.CheckAndSave(r.URL.Path))
}
