package adapterHTTP

import (
	"fmt"
	T "shortlink/internal/apptype"
	"time"

	"github.com/valyala/fasthttp"
)

var _ T.IHTTPClient = (*HTTPClientFast)(nil)

type HTTPClientFast struct {
	hcli *fasthttp.Client
	log  T.ILog
}

func NewHTTPClientFast(log T.ILog) *HTTPClientFast {
	return &HTTPClientFast{
		hcli: &fasthttp.Client{
			ReadTimeout:         10 * time.Second,
			WriteTimeout:        10 * time.Second,
			MaxIdleConnDuration: 10 * time.Second,
			MaxConnDuration:     10 * time.Second,
		},
		log: log,
	}
}

func (h HTTPClientFast) Get(link string) error {
	if code, _, err := h.hcli.Get(nil, link); err != nil {
		err = fmt.Errorf("%w: %w: %w", T.ErrHTTPClientFast, T.ErrGetMethod, err)
		h.log.LogError(err, "(HTTPClientFast).Get() http get method error")
		return err
	} else {
		if code < 500 {
			return nil
		} else {
			err := fmt.Errorf("%w: %w", T.ErrHTTPClientFast, T.ErrLinkNotAval)
			h.log.LogError(err, "(HTTPClientFast).Get() link is not available ")
			return err
		}
	}
}
