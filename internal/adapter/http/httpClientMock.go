package adapterHTTP

import (
	"net/http"
	T "shortlink/internal/apptype"
	"sync"
)

var _ T.IHTTPClient = (*HTTPClientMock)(nil)

type HTTPClientMock struct {
	hcli *sync.Map
}

func NewHTTPClientMock() HTTPClientMock {
	mockhcli := sync.Map{}
	mockhcli.Store("http://lib.ru", struct{}{})
	mockhcli.Store("http://lib.ru/PROZA/", struct{}{})
	mockhcli.Store("http://google.ru", struct{}{})
	return HTTPClientMock{
		hcli: &mockhcli,
	}
}

func (h *HTTPClientMock) Get(link string) (int, error) {
	_, ok := h.hcli.Load(link)
	if !ok {
		return http.StatusNotFound, nil
	}
	return http.StatusOK, nil
}
