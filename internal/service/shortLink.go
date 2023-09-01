package service

import (
	"net/http"
	"shortlink/internal/model"
	"sort"

	"shortlink/internal/i7e"
)

var _ i7e.ISvcShortLink = (*SvcShortLink)(nil)

type SvcShortLink struct {
	db   i7e.Idb
	hcli i7e.IHTTPClient
	log  i7e.ILog
}

func NewSvcShortLink(db i7e.Idb, hcli i7e.IHTTPClient, log i7e.ILog) SvcShortLink {
	return SvcShortLink{
		db:   db,
		hcli: hcli,
		log:  log,
	}
}

func (ssl *SvcShortLink) GetAllLinkPairs() []model.LinkPair {
	res := make([]model.LinkPair, 0, 8)
	allpairs := ssl.db.LoadAllLinkPairs()
	for _, el := range allpairs {
		if el.IsValid() {
			res = append(res, el)
		}
	}
	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Short() < res[j].Short()
	})
	return res
}

func (ssl *SvcShortLink) GetLinkLongFromLinkShort(linkshort string) model.LinkPair {
	lp := ssl.db.LoadLinkPair(linkshort)
	if lp.IsValid() {
		return lp
	}
	ssl.log.LogWarn("GetLinkLongFromLinkShort(): Link pair is not valid")
	return model.LinkPair{}
}

func (ssl *SvcShortLink) GetLinkShortFromLinkLong(linklong string) model.LinkPair {
	newLP := model.NewLinkPair(linklong)
	lp := ssl.db.LoadLinkPair(newLP.Short())
	if lp.IsValid() {
		return lp
	}
	ssl.log.LogWarn("GetLinkShortFromLinkLong(): Link pair is not valid")
	return model.LinkPair{}
}

func (ssl *SvcShortLink) SetLinkPairFromLinkLong(linklong string) model.LinkPair {
	empty := model.LinkPair{}
	newLP := model.NewLinkPair(linklong) // make pair
	if !newLP.IsValid() {
		ssl.log.LogWarn("SetLinkPairFromLinkLong(): Link pair is not valid")
		return empty
	}
	if !ssl.IsLinkLongHTTPValid(newLP.Long()) { // check http valid
		ssl.log.LogWarn("SetLinkPairFromLinkLong(): Link Long is not HTTP valid")
		return empty
	}
	dbsearchedLP := ssl.GetLinkLongFromLinkShort(newLP.Short()) // search in db
	if dbsearchedLP.IsValid() {
		return newLP
	}
	if !ssl.db.SaveLinkPair(newLP) { // save in db
		ssl.log.LogWarn("SetLinkPairFromLinkLong(): Link pair is not saved in db")
		return empty
	}
	return newLP
}

func (ssl *SvcShortLink) IsLinkLongHTTPValid(linklong string) bool {
	resp, err := ssl.hcli.Get(linklong)
	if err != nil {
		ssl.log.LogError(err, "IsLinkLongHttpValid(): http client GET error")
		return false
	}
	if resp == http.StatusOK {
		return true
	}
	ssl.log.LogWarn("IsLinkLongHttpValid(): http client GET is not OK: %d", resp)
	return false
}

/*
+check if long link is valid
+get the short link if you have a long link
+get the long link if you have a short link
+get ALL link pairs presented in db
*/
