package memstorage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemStorage(t *testing.T) {
	tests := []struct {
		name string
		key  string
		val  float64
	}{
		{
			name: "test1",
			key:  "metr1",
			val:  12332,
		},
	}
	Stor := NewStorage("1.json", false)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NoError(t, Stor.SetGauge(tt.key, tt.val))
		})
	}
}
