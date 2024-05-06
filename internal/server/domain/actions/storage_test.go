package actions

import (
	"context"
	"math/rand"
	"testing"

	"github.com/LobovVit/metric-collector/internal/server/config"
	"github.com/LobovVit/metric-collector/internal/server/domain/metrics"
)

func TestGetRepo(t *testing.T) {
	type args struct {
		ctx    context.Context
		config *config.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test get repo file", args: args{ctx: context.Background(), config: &config.Config{Restore: false, StoreInterval: 0, FileStoragePath: "/tmp/metrics-db.json"}}, wantErr: false},
		{name: "test get repo db", args: args{ctx: context.Background(), config: &config.Config{DSN: "postgresql://postgres", Restore: false, StoreInterval: 0, FileStoragePath: "/tmp/metrics-db.json"}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetRepo(tt.args.ctx, tt.args.config)
			if err != nil && !tt.wantErr {
				t.Errorf("GetRepo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepo(t *testing.T) {
	type args struct {
		ctx     context.Context
		config  *config.Config
		metrics []metrics.Metrics
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test file ok #1", args: args{ctx: context.Background(), config: &config.Config{Restore: false, StoreInterval: 0, FileStoragePath: "/tmp/metrics-db.json"},
			metrics: []metrics.Metrics{{ID: "k1", MType: "gauge"}, {ID: "k2", MType: "counter"}},
		}, wantErr: false},
		{name: "test db err", args: args{ctx: context.Background(), config: &config.Config{DSN: "postgresql://postgres", Restore: false, StoreInterval: 0, FileStoragePath: "/tmp/metrics-db.json"},
			metrics: []metrics.Metrics{{ID: "k1", MType: "gauge"}, {ID: "k2", MType: "counter"}},
		}, wantErr: true},
		{name: "test file ok #2", args: args{ctx: context.Background(),
			config:  &config.Config{Restore: true, StoreInterval: 0, FileStoragePath: "/tmp/metrics-db.json"},
			metrics: []metrics.Metrics{{ID: "k1", MType: "gauge"}, {ID: "k2", MType: "counter"}},
		}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := GetRepo(tt.args.ctx, tt.args.config)
			if err != nil && !tt.wantErr {
				t.Errorf("GetRepo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if r.storage != nil {
				err = r.storage.SetCounter(tt.args.ctx, "keyCounter", 111)
				if err != nil && !tt.wantErr {
					t.Errorf("SetCounter() error = %v, wantErr %v", err, tt.wantErr)
				}
				err = r.storage.SetGauge(tt.args.ctx, "keyGauge", 111.11)
				if err != nil && !tt.wantErr {
					t.Errorf("SetGauge() error = %v, wantErr %v", err, tt.wantErr)
				}
				valInt64 := rand.Int63n(100)
				valFloat64 := rand.Float64()
				tt.args.metrics[0].Value = &valFloat64
				tt.args.metrics[1].Delta = &valInt64
				err = r.SetBatch(tt.args.ctx, tt.args.metrics)
				if err != nil && !tt.wantErr {
					t.Errorf("SetBatch() error = %v, wantErr %v", err, tt.wantErr)
				}
				err = r.SaveToFile(tt.args.ctx)
				if err != nil && !tt.wantErr {
					t.Errorf("SaveToFile() error = %v, wantErr %v", err, tt.wantErr)
				}
				err = r.LoadFromFile(tt.args.ctx)
				if err != nil && !tt.wantErr {
					t.Errorf("LoadFromFile() error = %v, wantErr %v", err, tt.wantErr)
				}
				_, err := r.GetSingle(tt.args.ctx, tt.args.metrics[0].MType, "k1")
				if err != nil && !tt.wantErr {
					t.Errorf("GetSingleValText() error = %v, wantErr %v", err, tt.wantErr)
				}
				_, err = r.GetSingle(tt.args.ctx, tt.args.metrics[1].MType, "k2")
				if err != nil && !tt.wantErr {
					t.Errorf("GetSingleValText() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
