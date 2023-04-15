package services

import (
	"net/http"
	"shortlink/internal/models"
	"sort"

	"shortlink/internal/ports"
)

var _ ports.ISvcShortLink = (*SvcShortLink)(nil)

type SvcShortLink struct {
	db   ports.Idb
	hcli ports.IHttpClient
	log  ports.ILog
}

func NewSvcShortLink(db ports.Idb, hcli ports.IHttpClient, log ports.ILog) SvcShortLink {
	return SvcShortLink{
		db:   db,
		hcli: hcli,
		log:  log,
	}
}

func (ssl *SvcShortLink) GetAllLinkPairs() []models.LinkPair {
	res := make([]models.LinkPair, 0, 8)
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

func (ssl *SvcShortLink) GetLinkLongFromLinkShort(linkshort string) models.LinkPair {
	lp := ssl.db.LoadLinkPair(linkshort)
	if lp.IsValid() {
		return lp
	}
	ssl.log.LogWarn("GetLinkLongFromLinkShort(): Link pair is not valid")
	return models.LinkPair{}
}

func (ssl *SvcShortLink) GetLinkShortFromLinkLong(linklong string) models.LinkPair {
	newLP := models.NewLinkPair(linklong)
	lp := ssl.db.LoadLinkPair(newLP.Short())
	if lp.IsValid() {
		return lp
	}
	ssl.log.LogWarn("GetLinkShortFromLinkLong(): Link pair is not valid")
	return models.LinkPair{}
}

func (ssl *SvcShortLink) SetLinkPairFromLinkLong(linklong string) models.LinkPair {
	empty := models.LinkPair{}
	newLP := models.NewLinkPair(linklong) // make pair
	if !newLP.IsValid() {
		ssl.log.LogWarn("SetLinkPairFromLinkLong(): Link pair is not valid")
		return empty
	}
	if !ssl.IsLinkLongHttpValid(newLP.Long()) { // check http valid
		//ssl.log.LogWarn("SetLinkPairFromLinkLong(): Link Long is not HTTP valid")
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

func (ssl *SvcShortLink) IsLinkLongHttpValid(linklong string) bool {
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
