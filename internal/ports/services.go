package ports

import "shortlink/internal/models"

type IServShortLink interface {
	GetAllLinkPairs() []models.LinkPair
	GetLinkLongFromLinkShort(linkshort string) models.LinkPair
	GetLinkShortFromLinkLong(linklong string) models.LinkPair
	SetLinkPairFromLinkLong(linklong string) models.LinkPair
	IsLinkLongHttpValid(linklong string) bool
}
