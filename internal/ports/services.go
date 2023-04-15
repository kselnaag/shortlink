package ports

import "shortlink/internal/models"

type ISvcShortLink interface {
	GetAllLinkPairs() []models.LinkPair
	GetLinkLongFromLinkShort(linkshort string) models.LinkPair
	GetLinkShortFromLinkLong(linklong string) models.LinkPair
	SetLinkPairFromLinkLong(linklong string) models.LinkPair
	IsLinkLongHTTPValid(linklong string) bool
}
