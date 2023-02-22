package models

type Idb interface {
	SaveLinkPair(lp LinkPair) bool
	LoadLinkPair(ls string) LinkPair
	LoadAllLinkPairs() []LinkPair
}

/*


type App struct {
	//logger *Ilog
	//database *Idb
	//serviseSL *IserviceSL
	//httpClient *IHttpClient
	//httpServer   *IHttpServer
}

type Ilog interface {
	//
}

type IHttpServer interface {
	//
}
*/
