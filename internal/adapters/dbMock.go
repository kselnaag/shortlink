package adapters

import (
	"context"
	"shortlink/internal/models"
	"shortlink/internal/ports"
	"sync"
)

var _ ports.Idb = (*MockDB)(nil)

type MockDB struct {
	ctx *context.Context
	db  sync.Map
}

func NewMockDB(ctx *context.Context) MockDB {
	return MockDB{
		db:  sync.Map{},
		ctx: ctx,
	}
}

func (m *MockDB) SaveLinkPair(lp models.LinkPair) bool {
	m.db.Store(lp.Short, lp.Long)
	return true
}

func (m *MockDB) LoadLinkPair(linkshort string) models.LinkPair {
	res := models.LinkPair{}
	linklong, ok := m.db.Load(linkshort)
	if !ok {
		return res
	}
	res.Short = linkshort
	res.Long = linklong.(string)
	return res
}

func (m *MockDB) LoadAllLinkPairs() []models.LinkPair {
	res := make([]models.LinkPair, 0, 8)
	m.db.Range(func(key, value any) bool {
		res = append(res, models.LinkPair{Short: key.(string), Long: value.(string)})
		return true
	})
	return res
}
