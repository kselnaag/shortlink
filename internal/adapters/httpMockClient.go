package adapters

import (
	"net/http"
	"shortlink/internal/ports"
	"sync"
)

var _ ports.IHTTPClient = (*HTTPMockClient)(nil)

type HTTPMockClient struct {
	hcli *sync.Map
}

func NewHTTPMockClient() HTTPMockClient {
	mockhcli := sync.Map{}
	mockhcli.Store("http://lib.ru", struct{}{})
	mockhcli.Store("http://lib.ru/PROZA/", struct{}{})
	mockhcli.Store("http://google.ru", struct{}{})
	return HTTPMockClient{
		hcli: &mockhcli,
	}
}

func (h *HTTPMockClient) Get(link string) (int, error) {
	_, ok := h.hcli.Load(link)
	if !ok {
		return http.StatusNotFound, nil
	}
	return http.StatusOK, nil
}
