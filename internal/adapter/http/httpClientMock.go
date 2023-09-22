package adapterHTTP

import (
	"net/http"
	"shortlink/internal/types"
	"sync"
)

var _ types.IHTTPClient = (*HTTPClientMock)(nil)

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
