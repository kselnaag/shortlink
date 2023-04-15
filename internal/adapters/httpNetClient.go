package adapters

import (
	"net/http"
	"shortlink/internal/ports"
	"time"
)

var _ ports.IHTTPClient = (*HTTPNetClient)(nil)

type HTTPNetClient struct {
	hcli *http.Client
}

func NewHTTPNetClient() HTTPNetClient {
	transport := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}
	return HTTPNetClient{
		hcli: &http.Client{Transport: transport},
	}
}

func (h HTTPNetClient) Get(link string) (int, error) {
	resp, err := h.hcli.Get(link)
	if err != nil {
		return 0, err
	} else {
		defer resp.Body.Close()
	}
	return resp.StatusCode, nil
}
