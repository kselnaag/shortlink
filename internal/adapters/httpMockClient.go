package adapters

import (
	"net/http"
	"shortlink/internal/ports"
	"sync"
)

var _ ports.IHttpClient = (*HttpMockClient)(nil)

type HttpMockClient struct {
	hcli *sync.Map
}

func NewHttpMockClient() HttpMockClient {
	mockhcli := sync.Map{}
	mockhcli.Store("http://lib.ru", struct{}{})
	mockhcli.Store("http://lib.ru/PROZA/", struct{}{})
	mockhcli.Store("http://google.ru", struct{}{})
	return HttpMockClient{
		hcli: &mockhcli,
	}
}

func (h *HttpMockClient) Get(link string) (int, error) {
	_, ok := h.hcli.Load(link)
	if !ok {
		return http.StatusNotFound, nil
	}
	return http.StatusOK, nil
}
