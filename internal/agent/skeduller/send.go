package skeduller

import (
	"context"
	"fmt"
	"github.com/LobovVit/metric-collector/internal/agent/metrics"
	"github.com/go-resty/resty/v2"
	"strconv"
)

type restyClient struct {
	client *resty.Client
}

func (rc *restyClient) new() *restyClient {
	rc.client = resty.New()
	return rc
}

func sendRequest(ctx context.Context, metrics *metrics.Metrics, endPoint string) error {
	var sender restyClient
	sender.new()

	for k, v := range metrics.Gauge {
		_, err := sender.client.R().
			SetContext(ctx).
			SetHeader("Content-Type", "text/plain").
			Post(fmt.Sprintf("%vgauge/%v/%v", endPoint, k, strconv.FormatFloat(v, 'f', 10, 64)))
		if err != nil {
			return err
		}
	}
	for k, v := range metrics.Counter {
		_, err := sender.client.R().
			SetContext(ctx).
			SetHeader("Content-Type", "text/plain").
			Post(fmt.Sprintf("%vcounter/%v/%v", endPoint, k, strconv.FormatInt(v, 10)))
		if err != nil {
			return err
		}
	}
	return nil
}
