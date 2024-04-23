package metrics

import (
	"reflect"
	"testing"
)

func TestGetMetricStruct(t *testing.T) {
	tests := []struct {
		name string
		want *Metrics
	}{
		{name: "test struct",
			want: &Metrics{Metrics: make([]Metric, 32)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMetricStruct(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMetricStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetrics_GetMetrics(t *testing.T) {
	type fields struct {
		Metrics []Metric
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "test get metrics",
			fields: fields{Metrics: []Metric{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := GetMetricStruct()
			m.GetMetrics()
			if m.Metrics[1].ID != "Alloc" || m.Metrics[2].ID != "BuckHashSys" {
				t.Errorf("metrics didn't calculated %v ", m.Metrics[0].ID)
			}
		})
	}
}
