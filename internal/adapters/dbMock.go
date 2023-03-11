package adapters

import (
	"shortlink/internal/models"
	"shortlink/internal/ports"
	"sync"
)

var _ ports.Idb = (*MockDB)(nil)

type MockDB struct {
	db sync.Map
}

func NewMockDB() MockDB {
	return MockDB{
		db: sync.Map{},
	}
}

func (m *MockDB) SaveLinkPair(lp models.LinkPair) bool {
	m.db.Store(lp.Short(), lp.Long())
	return true
}

func (m *MockDB) LoadLinkPair(linkshort string) models.LinkPair {
	linklong, ok := m.db.Load(linkshort)
	if !ok {
		return models.LinkPair{}
	}
	return models.NewLinkPair(linklong.(string))
}

func (m *MockDB) LoadAllLinkPairs() []models.LinkPair {
	res := make([]models.LinkPair, 0, 8)
	m.db.Range(func(key, value any) bool {
		res = append(res, models.NewLinkPair(value.(string)))
		return true
	})
	return res
}
