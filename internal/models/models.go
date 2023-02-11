package models

type LinkPair struct {
	Short string
	Long  string
}

type Env struct {
	//db     *Idb
	//logger *Ilog
	//http   *Ihttp
}

type Idb interface {
	SaveLinkPair(lp LinkPair) bool
	LoadLinkPair(ls string) []LinkPair
	LoadAllLinkPairs() []LinkPair
}

type Ihttp interface {
	//
}

type Ilog interface {
	//
}
