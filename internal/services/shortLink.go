package services

import (
	"shortlink/internal/models"
	"shortlink/internal/ports"
)

var _ ports.IServShortLink = (*ServShortLink)(nil)

type ServShortLink struct {
	db   ports.Idb
	hcli ports.IHttpClient
}

func NewServShortLink(db ports.Idb, hcli ports.IHttpClient) ServShortLink {
	return ServShortLink{
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

func (ssl *ServShortLink) GetLinkShortFromLinkLong(linklong string) models.LinkPair {
	newLP := models.NewLinkPair(linklong)
	lp := ssl.db.LoadLinkPair(newLP.Short())
	if lp.IsValid() {
		return lp
	}
	return models.LinkPair{}
}

func (ssl *ServShortLink) SetLinkPairFromLinkLong(linklong string) models.LinkPair {
	empty := models.LinkPair{}
	newLP := models.NewLinkPair(linklong) // make pair
	if !newLP.IsValid() {
		return empty
	}
	if !ssl.IsLinkLongHttpValid(newLP.Long()) { // check http valid
		return empty
	}
	dbsearchedLP := ssl.GetLinkLongFromLinkShort(newLP.Short()) // search in db
	if dbsearchedLP.IsValid() {
		return newLP
	}
	if !ssl.db.SaveLinkPair(newLP) { // save in db
		return empty
	}
	return newLP
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
