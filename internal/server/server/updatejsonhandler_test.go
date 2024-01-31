package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LobovVit/metric-collector/internal/server/config"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type metric struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func testJSONRequest(t *testing.T, ts *httptest.Server, method, path string, met metric) (*http.Response, string) {

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

func TestUpdateJSONHandler(t *testing.T) {

	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name  string
		metod string
		path  string
		data  metric
		want  want
	}{
		{
			name:  "test1",
			metod: http.MethodPost,
			path:  "/update/",
			data:  metric{ID: "qqq", MType: "counter", Delta: int64(10)},
			want: want{
				contentType: "application/json",
				statusCode:  200,
			},
		},
		{
			name:  "test2",
			metod: http.MethodPost,
			path:  "/update/",
			data:  metric{ID: "qqq", MType: "gauge", Value: float64(10)},
			want: want{
				contentType: "application/json",
				statusCode:  200,
			},
		},
	}
	mux := chi.NewRouter()
	cfg, _ := config.GetConfig()
	tst := New(cfg)
	mux.Post("/update/", tst.updateJSONHandler)
	ts := httptest.NewServer(mux)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := testJSONRequest(t, ts, tt.metod, tt.path, tt.data)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			resp.Body.Close()
		})
	}
}
