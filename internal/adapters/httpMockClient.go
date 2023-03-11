package adapters

import (
	"shortlink/internal/ports"
	"sync"
)

var _ ports.IHttpClient = (*HttpMockClient)(nil)

type HttpMockClient struct {
	mock sync.Map
}

func NewHttpMockClient() HttpMockClient {
	mock := sync.Map{}
	mock.Store("http://lib.ru", struct{}{})
	return HttpMockClient{
		mock: mock,
	}
}

func (h HttpMockClient) Get(link string) (string, error) {
	_, ok := h.mock.Load(link)
	if !ok {
		return "404 Not Found", nil
	}
	return "200 OK", nil
}
