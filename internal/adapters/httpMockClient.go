package adapters

import (
	"context"
	"shortlink/internal/ports"
	"sync"
)

var _ ports.IHttpClient = (*HttpMockClient)(nil)

type HttpMockClient struct {
	ctx  *context.Context
	mock sync.Map
}

func NewHttpMockClient(ctx *context.Context) HttpMockClient {
	mock := sync.Map{}
	mock.Store("http://lib.ru", struct{}{})
	return HttpMockClient{
		ctx:  ctx,
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
