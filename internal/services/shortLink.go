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
	IsLinkLongHttpValid(linklong string) bool
}

type ServShortLink struct {
	ctx  *context.Context
	db   ports.Idb
	hcli ports.IHttpClient
}

func NewServShortLink(ctx *context.Context, db ports.Idb, hcli ports.IHttpClient) ServShortLink {
	return ServShortLink{
		ctx:  ctx,
		db:   db,
		hcli: hcli,
	}
}

func (ssl *ServShortLink) GetAllLinkPairs() []models.LinkPair {
	res := make([]models.LinkPair, 0, 8)
	allpairs := ssl.db.LoadAllLinkPairs()
	for _, el := range allpairs {
		if el.IsValid() {
			res = append(res, el)
		}
	}
	return res
}

func (ssl *ServShortLink) GetLinkLongFromLinkShort(linkshort string) models.LinkPair {
	lp := ssl.db.LoadLinkPair(linkshort)
	if lp.IsValid() {
		return lp
	}
	return models.LinkPair{}
}

func (ssl *ServShortLink) SetLinkPairFromLinkLong(linklong string) bool {
	newLP := models.NewLinkPair(linklong) // make pair
	if !newLP.IsValid() {
		return false
	}
	if !ssl.IsLinkLongHttpValid(newLP.Long()) { // check http valid
		return false
	}
	dbsearchedLP := ssl.GetLinkLongFromLinkShort(newLP.Short()) // search in db
	if dbsearchedLP.IsValid() {
		return true
	}
	if !ssl.db.SaveLinkPair(newLP) { // save in db
		return false
	}
	return true
}

func (ssl *ServShortLink) IsLinkLongHttpValid(linklong string) bool {
	resp, err := ssl.hcli.Get(linklong)
	if err != nil {
		return false
	}
	if resp == "200 OK" {
		return true
	}
	return false
}

/*
+check if long link is valid
+get the short link if you have a long link
+get the long link if you have a short link
+get ALL link pairs presented in db
*/
