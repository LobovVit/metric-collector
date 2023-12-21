package main

import (
	"github.com/LobovVit/metric-collector/internal/server/domain"
	"net/http"
	"strconv"
	"strings"
)

var storage = *domain.GetStorage()

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, updateHandler)
	return http.ListenAndServe(`:8080`, mux)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		//только POST
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(checkAndSave(r.URL.Path))
}

func checkAndSave(url string) int {
	part := strings.Split(url, `/`)
	if len(part) != 5 && len(part) != 6 {
		return http.StatusNotFound
	}
	switch part[2] {
	case "gauge":
		v, err := strconv.ParseFloat(part[4], 64)
		if err != nil {
			return http.StatusBadRequest
		}
		storage.SetGauge(part[3], v)
	case "counter":
		v, err := strconv.ParseInt(part[4], 10, 64)
		if err != nil {
			return http.StatusBadRequest
		}
		storage.SetCounter(part[3], v)
	default:
		return http.StatusBadRequest
	}
	//storage.GetAll()
	return http.StatusOK
}
