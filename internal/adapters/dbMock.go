package adapters

import (
	mod "shortlink/internal/models"
)

type mockDB struct {
	db map[string]string //  !!! sync + ctx
}

func NewMockDB() mockDB {
	return mockDB{
		db: make(map[string]string),
	}
}

func (mock *mockDB) SaveLinkPair(lp mod.LinkPair) bool {
	mock.db[lp.Short] = lp.Long
	return true
}

func (mock *mockDB) LoadLinkPair(ls string) []mod.LinkPair {
	result := make([]mod.LinkPair, 0, 1)
	ll, ok := mock.db[ls]
	if !ok {
		return result
	}
	return append(result, mod.LinkPair{Short: ls, Long: ll})
}

func (mock *mockDB) LoadAllLikPairs() []mod.LinkPair {
	result := make([]mod.LinkPair, 0, 8)
	for key, val := range mock.db {
		result = append(result, mod.LinkPair{Short: key, Long: val})
	}
	return result
}
