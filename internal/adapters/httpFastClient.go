package adapters

import (
	"shortlink/internal/ports"
	"time"

	"github.com/valyala/fasthttp"
)

var _ ports.IHttpClient = (*HttpFastClient)(nil)

type HttpFastClient struct {
	hcli *fasthttp.Client
}

func NewHttpFastClient() HttpFastClient {
	return HttpFastClient{
		hcli: &fasthttp.Client{
			ReadTimeout:         10 * time.Second,
			WriteTimeout:        10 * time.Second,
			MaxIdleConnDuration: 10 * time.Second,
			MaxConnDuration:     10 * time.Second,
		},
	}
}

func (h HttpFastClient) Get(link string) (int, error) {
	code, _, err := h.hcli.Get(nil, link)
	if err != nil {
		return 0, err
	}
	return code, nil
}
