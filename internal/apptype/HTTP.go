package types

type IHTTPClient interface {
	Get(link string) (int, error)
}

type IHTTPServer interface {
	Run() func(e error)
}

type HTTPMessage struct {
	IsResp bool
	Mode   string
	Body   string
}
