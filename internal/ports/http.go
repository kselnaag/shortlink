package ports

type IHttpClient interface {
	Get(ink string) (int, error)
}

type IHttpServer interface {
	Run() func()
}
