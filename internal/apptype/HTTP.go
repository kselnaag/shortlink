package types

type IHTTPClient interface {
	Get(link string) (int, error)
}

type IHTTPServer interface {
	Run() func(e error)
}
