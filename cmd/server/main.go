package main

import (
	"github.com/LobovVit/metric-collector/internal/server/handlers"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, handlers.UpdateHandler)
	return http.ListenAndServe(`:8080`, mux)
}
