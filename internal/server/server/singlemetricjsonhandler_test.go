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

func testJSONSingleRequest(t *testing.T, ts *httptest.Server, method, path string, met metric) (*http.Response, string) {

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

func TestUpdateJSONSingleHandler(t *testing.T) {

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
			path:  "/value/",
			data:  metric{ID: "qqq", MType: "counter"},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
			},
		},
		{
			name:  "test2",
			metod: http.MethodPost,
			path:  "/value/",
			data:  metric{ID: "www", MType: "gauge"},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
			},
		},
	}
	mux := chi.NewRouter()
	cfg, _ := config.GetConfig()
	tst := New(cfg)
	mux.Post("/value/", tst.singleMetricJSONHandler)
	ts := httptest.NewServer(mux)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := testJSONSingleRequest(t, ts, tt.metod, tt.path, tt.data)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			resp.Body.Close()
		})
	}
}
