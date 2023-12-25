package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestGetMetricStruct(t *testing.T) {
//	tests := []struct {
//		name string
//		want *Metrics
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := GetMetricStruct(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetMetricStruct() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestMetrics_GetMetrics(t *testing.T) {
	type fields struct {
		Gauge               map[string]float64
		Counter             map[string]int64
		CounterExecMemStats int64
	}
	tests := []struct {
		name string
		fil  fields
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fil: fields{
				Gauge:               map[string]float64{"metr1": 111},
				Counter:             map[string]int64{"metr1": 222},
				CounterExecMemStats: 333,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				Gauge:               tt.fil.Gauge,
				Counter:             tt.fil.Counter,
				CounterExecMemStats: tt.fil.CounterExecMemStats,
			}
			m.GetMetrics()
			assert.Equal(t, m.CounterExecMemStats, tt.fil.CounterExecMemStats+1)
		})
	}
}
