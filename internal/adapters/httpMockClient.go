package adapters

import (
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

func (h *HttpMockClient) Get(link string) (string, error) {
	_, ok := h.hcli.Load(link)
	if !ok {
		return "404 Not Found", nil
	}
	return "200 OK", nil
}
