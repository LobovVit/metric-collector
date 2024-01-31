package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type met struct {
	ID    string  `json:"id"`              // имя метрики
	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func TestCheckAndSave(t *testing.T) {
	tests := []struct {
		name string
		url  string
		tp   string
		code string
		val  string
		data met
	}{
		{
			name: "+ test #1",
			url:  "/update/counter/someMetric/527",
			tp:   "counter",
			code: "someMetric",
			val:  "527",
			data: met{ID: "someMetric", MType: "counter", Delta: int64(527)},
		}, {
			name: "+ test #2",
			url:  "/update/counter/someMetric/527/",
			tp:   "counter",
			code: "someMetric",
			val:  "527",
			data: met{ID: "someMetric", MType: "counter", Delta: int64(525)},
		}, {
			name: "- test #3",
			url:  "/update/counter/someMetric/527/ssa/",
			tp:   "counter",
			code: "someMetric",
			val:  "527",
			data: met{ID: "someMetric", MType: "counter", Delta: int64(27)},
		}, {
			name: "- test #4",
			url:  "/update/counter/some",
			tp:   "counter",
			code: "someMetric",
			val:  "527",
			data: met{ID: "someMetric", MType: "counter", Delta: int64(527)},
		}, {
			name: "- test #5",
			url:  "/update/counter/someMetric/hello",
			tp:   "counter",
			code: "someMetric",
			val:  "527",
			data: met{ID: "someMetric", MType: "counter", Delta: int64(527)},
		}, {
			name: "+ test #6",
			url:  "/update/gauge/someMetric/527",
			tp:   "gauge",
			code: "someMetric",
			val:  "527",
			data: met{ID: "someMetric", MType: "gauge", Value: float64(527)},
		}, {
			name: "+ test #7",
			url:  "/update/gauge/someMetric/527/",
			tp:   "gauge",
			code: "someMetric",
			val:  "527",
			data: met{ID: "someMetric", MType: "gauge", Value: float64(1)},
		}, {
			name: "- test #8",
			url:  "/update/gauge/someMetric/527/ssa/",
			tp:   "gauge",
			code: "someMetric",
			val:  "527",
			data: met{ID: "someMetric", MType: "gauge", Value: float64(1)},
		}, {
			name: "- test #9",
			url:  "/update/gauge/some",
			tp:   "gauge",
			code: "someMetric",
			val:  "444",
			data: met{ID: "someMetric", MType: "gauge", Value: float64(111)},
		}, {
			name: "- test #10",
			url:  "/update/gauge/someMetric/hello",
			tp:   "gauge",
			code: "someMetric",
			val:  "444",
			data: met{ID: "someMetric", MType: "counter", Delta: int64(527)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := GetRepo(false, 1000, "1.json")
			assert.NoError(t, x.CheckAndSaveText(tt.tp, tt.code, tt.val))
		})
	}
}
