package types

type IHTTPClient interface {
	Get(link string) error
}

type IHTTPServer interface {
	Run() func(e error)
}
