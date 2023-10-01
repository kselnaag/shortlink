package control

import (
	T "shortlink/internal/apptype"
	"shortlink/internal/model"
)

var _ T.ICtrlDB = (*CtrlDB)(nil)

type CtrlDB struct {
	db T.IDB
}

func NewCtrlDB(db T.IDB) CtrlDB {
	return CtrlDB{
		db: db,
	}
}

func (ctrl *CtrlDB) SaveLinkPair(lp model.LinkPair) bool {
	return ctrl.db.SaveLinkPair(T.DBlinksDTO{Short: lp.Short(), Long: lp.Long()})
}

func (ctrl *CtrlDB) LoadLinkPair(ls string) model.LinkPair {
	links := ctrl.db.LoadLinkPair(T.DBlinksDTO{Short: ls, Long: ""})
	return model.NewLinkPair(links.Long)
}

func (ctrl *CtrlDB) LoadAllLinkPairs() []model.LinkPair {
	res := make([]model.LinkPair, 0, 8)
	arrlinks := ctrl.db.LoadAllLinkPairs()
	for _, el := range arrlinks {
		res = append(res, model.NewLinkPair(el.Long))
	}
	return res
}
