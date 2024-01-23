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
		tp   string
		code string
		val  string
		want int
	}{
		{
			name: "+ test #1",
			url:  "/update/counter/someMetric/527",
			tp:   "counter",
			code: "someMetric",
			val:  "527",
			want: http.StatusOK,
		}, {
			name: "+ test #2",
			url:  "/update/counter/someMetric/527/",
			tp:   "counter",
			code: "someMetric",
			val:  "527",
			want: http.StatusOK,
		}, {
			name: "- test #3",
			url:  "/update/counter/someMetric/527/ssa/",
			tp:   "counter",
			code: "someMetric",
			val:  "527",
			want: http.StatusOK,
		}, {
			name: "- test #4",
			url:  "/update/counter/some",
			tp:   "counter",
			code: "someMetric",
			val:  "527",
			want: http.StatusOK,
		}, {
			name: "- test #5",
			url:  "/update/counter/someMetric/hello",
			tp:   "counter",
			code: "someMetric",
			val:  "527",
			want: http.StatusOK,
		}, {
			name: "+ test #6",
			url:  "/update/gauge/someMetric/527",
			tp:   "gauge",
			code: "someMetric",
			val:  "527",
			want: http.StatusOK,
		}, {
			name: "+ test #7",
			url:  "/update/gauge/someMetric/527/",
			tp:   "gauge",
			code: "someMetric",
			val:  "527",
			want: http.StatusOK,
		}, {
			name: "- test #8",
			url:  "/update/gauge/someMetric/527/ssa/",
			tp:   "gauge",
			code: "someMetric",
			val:  "527",
			want: http.StatusOK,
		}, {
			name: "- test #9",
			url:  "/update/gauge/some",
			tp:   "gauge",
			code: "someMetric",
			val:  "444",
			want: http.StatusOK,
		}, {
			name: "- test #10",
			url:  "/update/gauge/someMetric/hello",
			tp:   "gauge",
			code: "someMetric",
			val:  "444",
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := GetRepo()
			assert.NoError(t, x.CheckAndSaveText(tt.tp, tt.code, tt.val), tt.want)
		})
	}
}
