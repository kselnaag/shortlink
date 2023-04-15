package adapters

import (
	"shortlink/internal/models"
	"shortlink/internal/ports"
	"sync"
)

var _ ports.Idb = (*DBMock)(nil)

type DBMock struct {
	cfg *CfgEnv
	db  *sync.Map
}

func NewDBMock(cfg *CfgEnv) DBMock {
	dbmock := sync.Map{}
	dbmock.Store("5clp60", "http://lib.ru")
	dbmock.Store("dhiu79", "http://google.ru")
	return DBMock{
		cfg: cfg,
		db:  &dbmock,
	}
}

func (m *DBMock) SaveLinkPair(lp models.LinkPair) bool {
	m.db.Store(lp.Short(), lp.Long())
	return true
}

func (m *DBMock) LoadLinkPair(linkshort string) models.LinkPair {
	linklong, ok := m.db.Load(linkshort)
	if !ok {
		return models.LinkPair{}
	}
	return models.NewLinkPair(linklong.(string))
}

func (m *DBMock) LoadAllLinkPairs() []models.LinkPair {
	res := make([]models.LinkPair, 0, 8)
	m.db.Range(func(key, value any) bool {
		res = append(res, models.NewLinkPair(value.(string)))
		return true
	})
	return res
}
