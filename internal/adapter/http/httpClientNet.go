package adapterHTTP

import (
	"fmt"
	"net/http"
	T "shortlink/internal/apptype"
	"time"
)

var _ T.IHTTPClient = (*HTTPClientNet)(nil)

type HTTPClientNet struct {
	hcli *http.Client
	log  T.ILog
}

func NewHTTPClientNet(log T.ILog) *HTTPClientNet {
	transport := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}
	return &HTTPClientNet{
		hcli: &http.Client{Transport: transport},
		log:  log,
	}
}

func (h HTTPClientNet) Get(link string) error {
	if resp, err := h.hcli.Get(link); err != nil { //nolint:noctx // do it simple
		err = fmt.Errorf("%w: %w: %w", T.ErrHTTPClientNet, T.ErrGetMethod, err)
		h.log.LogError(err, "(HTTPClientNet).Get() http get method error")
		return err
	} else {
		defer func() {
			_ = resp.Body.Close()
		}()
		if resp.StatusCode < 500 {
			return nil
		} else {
			err := fmt.Errorf("%w: %w", T.ErrHTTPClientNet, T.ErrLinkNotAval)
			h.log.LogError(err, "(HTTPClientNet).Get(): link is not available")
			return err
		}
	}
}
