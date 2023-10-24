package adapterHTTP

import (
	"errors"
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
	if resp, err := h.hcli.Get(link); err != nil {
		h.log.LogError(err, "(HTTPClientNet).Get() http error ")
		return err
	} else {
		defer resp.Body.Close()
		if resp.StatusCode < 500 {
			return nil
		} else {
			err := errors.New("(HTTPClientNet).Get(): Link is not available")
			h.log.LogError(err, "(HTTPClientNet).Get() http error ")
			return err
		}
	}
}
