package types

import "shortlink/internal/model"

type HTTPMessageDTO struct {
	IsResp bool   `json:"IsResp"` //nolint:tagliatelle // wrong json camel
	Mode   string `json:"Mode"`   //nolint:tagliatelle
	Body   string `json:"Body"`   //nolint:tagliatelle
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
