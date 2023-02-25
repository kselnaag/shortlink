package adapters

import (
	"context"
	"net/http"
	"shortlink/internal/ports"
	"time"
)

var _ ports.IHttpClient = (*HttpNetClient)(nil)

type HttpNetClient struct {
	ctx  *context.Context
	hcli *http.Client
}

func NewHttpNetClient(ctx *context.Context) HttpNetClient {
	transport := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}
	return HttpNetClient{
		ctx:  ctx,
		hcli: &http.Client{Transport: transport},
	}
}

func (h HttpNetClient) Get(link string) (string, error) {
	resp, err := h.hcli.Get(link)
	if err != nil {
		return "", err
	}
	return resp.Status, nil
}
