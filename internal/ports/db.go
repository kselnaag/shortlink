package ports

import (
	"shortlink/internal/models"
)

type Idb interface {
	SaveLinkPair(lp models.LinkPair) bool
	LoadLinkPair(ls string) models.LinkPair
	LoadAllLinkPairs() []models.LinkPair
}
