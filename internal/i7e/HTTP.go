package i7e

type IHTTPClient interface {
	Get(ink string) (int, error)
}

type IHTTPServer interface {
	Run() func()
}
