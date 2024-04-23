package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LobovVit/metric-collector/internal/server/config"
)

var TS = ServerRun()

type metric struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func ServerRun() *httptest.Server {
	mux := chi.NewRouter()
	cfg := config.Config{Host: "localhost:8080", LogLevel: "info", StoreInterval: 100, FileStoragePath: "1.json", Restore: false}
	tst, _ := New(context.Background(), &cfg)
	mux.Get("/", tst.allMetricsHandler)
	mux.Get("/ping", tst.dbPingHandler)
	mux.Post("/value/", tst.singleMetricJSONHandler)
	mux.Get("/value/{type}/{name}", tst.singleMetricHandler)
	mux.Post("/update/", tst.updateJSONHandler)
	mux.Post("/update/{type}/{name}/{value}", tst.updateHandler)
	return httptest.NewServer(mux)
}

func request(t *testing.T, ts *httptest.Server, method, path string, met metric) (*http.Response, string) {

	data, err := json.Marshal(&met)
	require.NoError(t, err)
	b := bytes.NewBuffer(data)
	req, err := http.NewRequest(method, ts.URL+path, b)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestHandlers(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		method1 string
		method2 string
		path1   string
		path2   string
		path3   string
		path4   string
		path5   string
		data1   metric
		data2   metric
		want    want
	}{
		{name: "test handlers #1",
			method1: http.MethodPost,
			method2: http.MethodGet,
			path1:   "/value/",
			path2:   "/update/",
			path3:   "/update/counter/someMetric/527",
			path4:   "/value/counter/someMetric",
			path5:   "/ping",
			data1:   metric{ID: "www", MType: "gauge", Value: float64(10)},
			data2:   metric{},
			want: want{
				contentType: "application/json",
				statusCode:  200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updResp, updBody := request(t, TS, tt.method1, tt.path2, tt.data1)
			assert.Equal(t, tt.want.statusCode, updResp.StatusCode)
			updResp.Body.Close()
			getResp, getJBody := request(t, TS, tt.method1, tt.path1, tt.data1)
			assert.Equal(t, tt.want.statusCode, getResp.StatusCode)
			assert.Equal(t, tt.want.contentType, getResp.Header.Get("Content-Type"))
			getResp.Body.Close()
			assert.Equal(t, updBody, getJBody)
			updTextResp, _ := request(t, TS, tt.method1, tt.path3, tt.data2)
			assert.Equal(t, tt.want.statusCode, updTextResp.StatusCode)
			updTextResp.Body.Close()
			getTextResp, _ := request(t, TS, tt.method2, tt.path4, tt.data2)
			assert.Equal(t, tt.want.statusCode, getTextResp.StatusCode)
			getTextResp.Body.Close()
			getPingResp, _ := request(t, TS, tt.method2, tt.path5, tt.data2)
			assert.Equal(t, http.StatusInternalServerError, getPingResp.StatusCode)
			getPingResp.Body.Close()
		})
	}
}

func ExampleServer_Run() {
	cfg := config.Config{Host: "localhost:8080", LogLevel: "info", StoreInterval: 100, FileStoragePath: "1.json", Restore: false}
	example, _ := New(context.Background(), &cfg)

	mux := chi.NewRouter()
	mux.Get("/", example.allMetricsHandler)
	mux.Get("/ping", example.dbPingHandler)
	mux.Post("/value/", example.singleMetricJSONHandler)
	mux.Get("/value/{type}/{name}", example.singleMetricHandler)
	mux.Post("/update/", example.updateJSONHandler)
	mux.Post("/update/{type}/{name}/{value}", example.updateHandler)
	httpServer := &http.Server{
		Addr:    cfg.Host,
		Handler: mux,
	}
	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
