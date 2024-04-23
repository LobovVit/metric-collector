package memstorage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage(t *testing.T) {
	tests := []struct {
		name string
		key  string
		val  float64
	}{
		{
			name: "test nem storage",
			key:  "metric1",
			val:  12332,
		},
	}
	Stor, _ := NewStorage(context.Background(), false, 100, "1.json")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, Stor.SetGauge(context.Background(), tt.key, tt.val))
		})
	}
}
