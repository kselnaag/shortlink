package types

import "shortlink/internal/model"

type Idb interface {
	SaveLinkPair(lp model.LinkPair) bool
	LoadLinkPair(ls string) model.LinkPair
	LoadAllLinkPairs() []model.LinkPair
}
