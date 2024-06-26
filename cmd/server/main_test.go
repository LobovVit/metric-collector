package main

import (
	"context"
	"testing"
	"time"
)

func Test_run(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "check run server",
			wantErr: true},
	}
	for _, tt := range tests {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		t.Run(tt.name, func(t *testing.T) {
			if err := run(ctx); (err != nil) != tt.wantErr {
				cancel()
			}
		})
	}
}
