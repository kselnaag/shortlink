package ports

type IHTTPClient interface {
	Get(ink string) (int, error)
}

type IHTTPServer interface {
	Run() func()
}
