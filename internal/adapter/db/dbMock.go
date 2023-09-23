package adapterDB

import (
	T "shortlink/internal/apptype"
	"shortlink/internal/model"
	"sync"
)

var _ T.Idb = (*DBMock)(nil)

type DBMock struct {
	cfg *T.CfgEnv
	db  *sync.Map
}

func NewDBMock(cfg *T.CfgEnv) DBMock {
	dbmock := sync.Map{}
	dbmock.Store("5clp60", "http://lib.ru")
	dbmock.Store("dhiu79", "http://google.ru")
	return DBMock{
		cfg: cfg,
		db:  &dbmock,
	}
}

func (m *DBMock) SaveLinkPair(lp model.LinkPair) bool {
	m.db.Store(lp.Short(), lp.Long())
	return true
}

func (m *DBMock) LoadLinkPair(linkshort string) model.LinkPair {
	linklong, ok := m.db.Load(linkshort)
	if !ok {
		return model.LinkPair{}
	}
	return model.NewLinkPair(linklong.(string))
}

func (m *DBMock) LoadAllLinkPairs() []model.LinkPair {
	res := make([]model.LinkPair, 0, 8)
	m.db.Range(func(key, value any) bool {
		res = append(res, model.NewLinkPair(value.(string)))
		return true
	})
	return res
}
