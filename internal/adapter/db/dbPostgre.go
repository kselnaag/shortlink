package adapterDB

import (
	T "shortlink/internal/apptype"
)

var _ T.IDB = (*DBPostgre)(nil)

type DBPostgre struct {
	cfg *T.CfgEnv
	log T.ILog
}

func NewDBPostgre(cfg *T.CfgEnv, log T.ILog) DBPostgre {
	return DBPostgre{
		cfg: cfg,
		log: log,
	}
}

func (p *DBPostgre) SaveLinkPair(links T.DBlinksDTO) bool {
	/* m.db.Store(lp.Short(), lp.Long()) */
	return true
}

func (p *DBPostgre) LoadLinkPair(links T.DBlinksDTO) T.DBlinksDTO { // linkshort
	/* linklong, ok := m.db.Load(linkshort)
	if !ok {
		return model.LinkPair{}
	} */
	return T.DBlinksDTO{}
}

func (p *DBPostgre) LoadAllLinkPairs() []T.DBlinksDTO {
	/* res := make([]model.LinkPair, 0, 8)
	m.db.Range(func(key, value any) bool {
		res = append(res, model.NewLinkPair(value.(string)))
		return true
	}) */
	return []T.DBlinksDTO{}
}

func (p *DBPostgre) Connect() func(e error) {

	return func(e error) {
		if e != nil {
			p.log.LogError(e, "DBPostgre.Connect(): db graceful_shutdown error")
		}
	}
}
