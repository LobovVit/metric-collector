package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LobovVit/metric-collector/internal/server/config"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestUpdateHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name  string
		metod string
		path  string
		want  want
	}{
		{
			name:  "test1",
			metod: http.MethodPost,
			path:  "/update/counter/someMetric/527",
			want: want{
				contentType: "text/plain",
				statusCode:  200,
			},
		},
	}
	mux := chi.NewRouter()
	cfg, _ := config.GetConfig()
	tst := New(cfg)
	mux.Post("/update/{type}/{name}/{value}", tst.updateHandler)
	ts := httptest.NewServer(mux)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, tt.metod, tt.path)
			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			resp.Body.Close()
		})
	}
}
