package types

import "shortlink/internal/model"

type ICtrlDB interface {
	SaveLinkPair(lp model.LinkPair) bool
	LoadLinkPair(ls string) model.LinkPair
	LoadAllLinkPairs() []model.LinkPair
}
