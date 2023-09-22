package service

import (
	"net/http"
	"shortlink/internal/model"
	"sort"

	"shortlink/internal/types"
)

var _ types.ISvcShortLink = (*SvcShortLink)(nil)

type SvcShortLink struct {
	db   types.Idb
	hcli types.IHTTPClient
	log  types.ILog
}

func NewSvcShortLink(db types.Idb, hcli types.IHTTPClient, log types.ILog) SvcShortLink {
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
	ssl.log.LogDebug("GetLinkLongFromLinkShort(): Link pair is not valid")
	return model.LinkPair{}
}

func (ssl *SvcShortLink) GetLinkShortFromLinkLong(linklong string) model.LinkPair {
	newLP := model.NewLinkPair(linklong)
	lp := ssl.db.LoadLinkPair(newLP.Short())
	if lp.IsValid() {
		return lp
	}
	ssl.log.LogDebug("GetLinkShortFromLinkLong(): Link pair is not valid")
	return model.LinkPair{}
}

func (ssl *SvcShortLink) SetLinkPairFromLinkLong(linklong string) model.LinkPair {
	empty := model.LinkPair{}
	newLP := model.NewLinkPair(linklong) // make pair
	if !newLP.IsValid() {
		ssl.log.LogDebug("SetLinkPairFromLinkLong(): Link pair is not valid")
		return empty
	}
	if !ssl.IsLinkLongHTTPValid(newLP.Long()) { // check http valid
		ssl.log.LogDebug("SetLinkPairFromLinkLong(): Link Long is not HTTP valid")
		return empty
	}
	dbsearchedLP := ssl.GetLinkLongFromLinkShort(newLP.Short()) // search in db
	if dbsearchedLP.IsValid() {
		return newLP
	}
	if !ssl.db.SaveLinkPair(newLP) { // save in db
		ssl.log.LogDebug("SetLinkPairFromLinkLong(): Link pair is not saved in db")
		return empty
	}
	return newLP
}

func (ssl *SvcShortLink) IsLinkLongHTTPValid(linklong string) bool {
	resp, err := ssl.hcli.Get(linklong)
	if err != nil {
		ssl.log.LogDebug("IsLinkLongHttpValid(): http client GET error: %s", err.Error())
		return false
	}
	if resp == http.StatusOK {
		return true
	}
	ssl.log.LogDebug("IsLinkLongHttpValid(): http client GET is not OK: %d", resp)
	return false
}

/*
+check if long link is valid
+get the short link if you have a long link
+get the long link if you have a short link
+get ALL link pairs presented in db
*/
