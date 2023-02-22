package services

import (
	"context"
	"shortlink/internal/models"
	"shortlink/internal/ports"
)

var _ IServShortLink = (*ServShortLink)(nil)

type IServShortLink interface {
	GetAllLinkPairs() []models.LinkPair
	GetLinkLongFromLinkShort(linkshort string) models.LinkPair
	SetLinkPairFromLinkLong(linklong string) bool
}

type ServShortLink struct {
	db         *ports.Idb
	httpClient *IHttpClient
	ctx        *context.Context
	// log        *Ilog
	// httpServer *IHttpServer
}

/* func NewServiceSL(ctx, db, cli, serv, log) {
	return ServiceSL{
		db
	}
} */

func (ssl *ServShortLink) GetAllLinkPairs() []models.LinkPair {
	result := make([]models.LinkPair, 0, 8)
	allpairs := (*ssl.db).LoadAllLinkPairs()
	for _, el := range allpairs {
		if el.IsValid() {
			result = append(result, el)
		}
	}
	return result
}

func (ssl *ServShortLink) GetLinkLongFromLinkShort(linkshort string) models.LinkPair {
	lp := (*ssl.db).LoadLinkPair(linkshort)
	if lp.IsValid() {
		return lp
	}
	return models.LinkPair{}
}

func (ssl *ServShortLink) SetLinkPairFromLinkLong(linklong string) bool {
	// make pair
	// search in db
	// check http valid
	// save in db

	newLP := models.NewLinkPair(linklong)
	res := newLP.IsValid()

	dbsearchedLP := ssl.GetLinkLongFromLinkShort(newLP.Short)
	res = dbsearchedLP.IsValid()

	res = (*ssl.HttpClient).IsLinkLongHttpValid(newLP.Long)

	res = (*ssl.db).SaveLinkPair(newLP)

	return res
}

func (ssl *ServShortLink) isLinkLongHttpValid(linklong string) bool {
	// HTTP !
	return true
}

/*
redirect from short link to long link
send html UI

-check if long link is valid
+get the short link if you have a long link
+get the long link if you have a short link
+get ALL link pairs presented in db
*/
