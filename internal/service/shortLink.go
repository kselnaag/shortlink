package service

import (
	"shortlink/internal/model"
	"sort"

	T "shortlink/internal/apptype"
)

var _ T.ISvcShortLink = (*SvcShortLink)(nil)

type SvcShortLink struct {
	ctrl T.ICtrlDB
	hcli T.IHTTPClient
	log  T.ILog
}

func NewSvcShortLink(ctrl T.ICtrlDB, hcli T.IHTTPClient, log T.ILog) *SvcShortLink {
	return &SvcShortLink{
		ctrl: ctrl,
		hcli: hcli,
		log:  log,
	}
}

func (ssl *SvcShortLink) GetAllLinkPairs() []model.LinkPair {
	res := make([]model.LinkPair, 0, 8)
	allpairs := ssl.ctrl.LoadAllLinkPairs()
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
	lp := ssl.ctrl.LoadLinkPair(linkshort)
	if lp.IsValid() {
		return lp
	}
	ssl.log.LogDebug("GetLinkLongFromLinkShort(): Link pair is not valid")
	return model.LinkPair{}
}

func (ssl *SvcShortLink) GetLinkShortFromLinkLong(linklong string) model.LinkPair {
	newLP := model.NewLinkPair(linklong)
	lp := ssl.ctrl.LoadLinkPair(newLP.Short())
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
	if !ssl.ctrl.SaveLinkPair(newLP) { // save in db
		ssl.log.LogDebug("SetLinkPairFromLinkLong(): Link pair is not saved in db")
		return empty
	}
	return newLP
}

func (ssl *SvcShortLink) IsLinkLongHTTPValid(linklong string) bool {
	if err := ssl.hcli.Get(linklong); err != nil {
		ssl.log.LogDebug("IsLinkLongHttpValid(): http client GET error: %s", err.Error())
		return false
	}
	return true
}
