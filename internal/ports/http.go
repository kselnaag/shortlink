package ports

type IHttpClient interface {
	Get(ink string) (string, error)
}

type IHttpServer interface {
	//
}
