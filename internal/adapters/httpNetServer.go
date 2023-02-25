package adapters

import (
	"context"
	"net/http"
	"shortlink/internal/ports"
)

var _ ports.IHttpServer = (*HttpNetServer)(nil)

type HttpNetServer struct {
	ctx  *context.Context
	hcli *http.Server
}

/*
redirect from short link to long link
send html UI

(check if long link is valid)
(get the short link if you have a long link)
(get the long link if you have a short link)
(get ALL link pairs presented in db)
*/
