package adapterHTTP

import (
	"errors"
	T "shortlink/internal/apptype"
	"sync"
)

var _ T.IHTTPClient = (*HTTPClientMock)(nil)

type HTTPClientMock struct {
	hcli *sync.Map
	log  T.ILog
}

func NewHTTPClientMock(log T.ILog) *HTTPClientMock {
	mockhcli := sync.Map{}
	mockhcli.Store("http://lib.ru", struct{}{})
	mockhcli.Store("http://lib.ru/PROZA/", struct{}{})
	mockhcli.Store("http://google.ru", struct{}{})
	return &HTTPClientMock{
		hcli: &mockhcli,
		log:  log,
	}
}

func (h *HTTPClientMock) Get(link string) error {
	_, ok := h.hcli.Load(link)
	if !ok {
		return errors.New("(HTTPClientMockt).Get(): Link is not available")
	}
	return nil
}
