package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.metod, tt.path, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(UpdateHandler)
			h(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			defer result.Body.Close()
		})
	}
}
