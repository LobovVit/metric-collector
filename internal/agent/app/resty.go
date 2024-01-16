package app

import (
	"github.com/go-resty/resty/v2"
)

type restyClient struct {
	client *resty.Client
}

func (rc *restyClient) new() *restyClient {
	rc.client = resty.New()
	return rc
}
