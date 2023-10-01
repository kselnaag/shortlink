package types

type DBlinksDTO struct {
	Long  string
	Short string
}

type IDB interface {
	SaveLinkPair(links DBlinksDTO) bool
	LoadLinkPair(links DBlinksDTO) DBlinksDTO
	LoadAllLinkPairs() []DBlinksDTO
	Connect() func(e error)
}
