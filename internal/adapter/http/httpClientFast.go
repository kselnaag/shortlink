package adapterHTTP

import (
	I "shortlink/internal/intrface"
	"time"

	"github.com/valyala/fasthttp"
)

var _ I.IHTTPClient = (*HTTPClientFast)(nil)

type HTTPClientFast struct {
	hcli *fasthttp.Client
}

func NewHTTPClientFast() HTTPClientFast {
	return HTTPClientFast{
		hcli: &fasthttp.Client{
			ReadTimeout:         10 * time.Second,
			WriteTimeout:        10 * time.Second,
			MaxIdleConnDuration: 10 * time.Second,
			MaxConnDuration:     10 * time.Second,
		},
	}
}

func (h HTTPClientFast) Get(link string) (int, error) {
	code, _, err := h.hcli.Get(nil, link)
	if err != nil {
		return 0, err
	}
	return code, nil
}
