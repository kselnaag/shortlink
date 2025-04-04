package adapterDB

import (
	T "shortlink/internal/apptype"
	"sync"
)

var _ T.IDB = (*DBMock)(nil)

type DBMock struct {
	log T.ILog
	cfg *T.CfgEnv
	db  *sync.Map
}

func NewDBMock(cfg *T.CfgEnv, log T.ILog) *DBMock {
	dbmock := sync.Map{}
	dbmock.Store("5clp60", "http://lib.ru")
	dbmock.Store("dhiu79", "http://google.ru")
	return &DBMock{
		log: log,
		cfg: cfg,
		db:  &dbmock,
	}
}

func (m *DBMock) SaveLinkPair(links T.DBlinksDTO) bool {
	m.db.Store(links.Short, links.Long)
	return true
}

func (m *DBMock) LoadLinkPair(links T.DBlinksDTO) T.DBlinksDTO {
	linklong, ok := m.db.Load(links.Short)
	if !ok {
		return T.DBlinksDTO{}
	}
	linklongS, _ := linklong.(string)
	return T.DBlinksDTO{Short: links.Short, Long: linklongS}
}

func (m *DBMock) LoadAllLinkPairs() []T.DBlinksDTO {
	res := make([]T.DBlinksDTO, 0, 8)
	m.db.Range(func(key, value any) bool {
		keyS, _ := key.(string)
		valS, _ := value.(string)
		res = append(res, T.DBlinksDTO{Short: keyS, Long: valS})
		return true
	})
	return res
}

func (m *DBMock) Connect() func(e error) {
	m.log.LogInfo("mock db connected")
	return func(e error) {
		if e != nil {
			m.log.LogError(e, "DBMock.Connect(): db graceful_shutdown error")
		}
		m.log.LogInfo("mock db disconnected")
	}
}
