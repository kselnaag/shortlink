package adapters

import (
	"net/http"
	"shortlink/internal/ports"
	"time"
)

var _ ports.IHttpClient = (*HttpNetClient)(nil)

type HttpNetClient struct {
	hcli *http.Client
}

func NewHttpNetClient() HttpNetClient {
	transport := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}
	return HttpNetClient{
		hcli: &http.Client{Transport: transport},
	}
}

func (h HttpNetClient) Get(link string) (int, error) {
	resp, err := h.hcli.Get(link)
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}
