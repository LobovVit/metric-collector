package actions

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCheckAndSave(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "+ test #1",
			url:  "/update/counter/someMetric/527",
			want: http.StatusOK,
		}, {
			name: "+ test #2",
			url:  "/update/counter/someMetric/527/",
			want: http.StatusOK,
		}, {
			name: "- test #3",
			url:  "/update/counter/someMetric/527/ssa/",
			want: http.StatusNotFound,
		}, {
			name: "- test #4",
			url:  "/update/counter/some",
			want: http.StatusNotFound,
		}, {
			name: "- test #5",
			url:  "/update/counter/someMetric/hello",
			want: http.StatusBadRequest,
		}, {
			name: "+ test #6",
			url:  "/update/gauge/someMetric/527",
			want: http.StatusOK,
		}, {
			name: "+ test #7",
			url:  "/update/gauge/someMetric/527/",
			want: http.StatusOK,
		}, {
			name: "- test #8",
			url:  "/update/gauge/someMetric/527/ssa/",
			want: http.StatusNotFound,
		}, {
			name: "- test #9",
			url:  "/update/gauge/some",
			want: http.StatusNotFound,
		}, {
			name: "- test #10",
			url:  "/update/gauge/someMetric/hello",
			want: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, CheckAndSave(tt.url), tt.want)
		})
	}
}
