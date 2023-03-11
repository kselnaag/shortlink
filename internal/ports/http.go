package ports

import (
	"net/http"
)

type IHttpClient interface {
	Get(ink string) (string, error)
}

type IHttpServer interface {
	Handle() IHttpServer
	Run(port string) *http.Server
}
