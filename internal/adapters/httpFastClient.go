package adapters

import (
	"shortlink/internal/ports"
	"time"

	"github.com/valyala/fasthttp"
)

var _ ports.IHTTPClient = (*HTTPFastClient)(nil)

type HTTPFastClient struct {
	hcli *fasthttp.Client
}

func NewHttpFastClient() HTTPFastClient {
	return HTTPFastClient{
		hcli: &fasthttp.Client{
			ReadTimeout:         10 * time.Second,
			WriteTimeout:        10 * time.Second,
			MaxIdleConnDuration: 10 * time.Second,
			MaxConnDuration:     10 * time.Second,
		},
	}
}

func (h HTTPFastClient) Get(link string) (int, error) {
	code, _, err := h.hcli.Get(nil, link)
	if err != nil {
		return 0, err
	}
	return code, nil
}
