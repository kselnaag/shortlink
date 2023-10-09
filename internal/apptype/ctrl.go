package types

import "shortlink/internal/model"

type HTTPMessageDTO struct {
	IsResp bool
	Mode   string
	Body   string
}

type ICtrlHTTP interface {
	AllPairs() (string, error)
	Long(body []byte) (string, error)
	Short(body []byte) (string, error)
	Save(body []byte) (string, error)
	Hash(hash string) (string, error)
}

type ICtrlDB interface {
	SaveLinkPair(lp model.LinkPair) bool
	LoadLinkPair(ls string) model.LinkPair
	LoadAllLinkPairs() []model.LinkPair
}
