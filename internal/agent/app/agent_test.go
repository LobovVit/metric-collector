package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/LobovVit/metric-collector/internal/agent/config"
	"github.com/LobovVit/metric-collector/internal/agent/metrics"
)

func TestAgent_Run(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "test1", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.GetConfig()
			if err != nil {
				t.Fatal(err)
			}
			a := New(cfg)
			ctx, chancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer chancel()
			if err := a.Run(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAgent_sendRequest(t *testing.T) {
	type fields struct {
		cfg    *config.Config
		client *resty.Client
	}
	type args struct {
		ctx     context.Context
		metrics *metrics.Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "test batch mode", fields: fields{cfg: &config.Config{ReportFormat: "batch", MaxCntInBatch: 10, RateLimit: 3}, client: resty.New()}, args: args{ctx: context.Background(), metrics: metrics.GetMetricStruct()}, wantErr: false},
		{name: "test json mode", fields: fields{cfg: &config.Config{ReportFormat: "json", MaxCntInBatch: 10, RateLimit: 3}, client: resty.New()}, args: args{ctx: context.Background(), metrics: metrics.GetMetricStruct()}, wantErr: false},
		{name: "test text mode", fields: fields{cfg: &config.Config{ReportFormat: "text", MaxCntInBatch: 10, RateLimit: 3}, client: resty.New()}, args: args{ctx: context.Background(), metrics: metrics.GetMetricStruct()}, wantErr: false},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	defer srv.Close()
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			tt.fields.cfg.Host = srv.URL + "/"
			tt.fields.client = resty.New()
			tt.args.metrics.GetMetrics()
			a := Agent(tt.fields)
			if err := a.sendRequest(tt.args.ctx, tt.args.metrics); (err != nil) != tt.wantErr {
				t.Errorf("sendRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
