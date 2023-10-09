package adapterDB

import (
	"context"
	T "shortlink/internal/apptype"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ T.IDB = (*DBPostgre)(nil)

type DBPostgre struct {
	cfg  *T.CfgEnv
	log  T.ILog
	conn *pgxpool.Pool
}

func NewDBPostgre(cfg *T.CfgEnv, log T.ILog) DBPostgre {
	return DBPostgre{
		cfg: cfg,
		log: log,
	}
}

func (p *DBPostgre) SaveLinkPair(links T.DBlinksDTO) bool {
	ctx := context.Background()
	query := "INSERT INTO shortlink VALUES ($1, $2)"
	tag, err := p.conn.Exec(ctx, query, links.Short, links.Long)
	if err != nil {
		p.log.LogError(err, "(DBPostgre).SaveLinkPair(): INSERT error")
		return false
	}
	p.log.LogDebug("(DBPostgre).SaveLinkPair(): %s", tag)
	return true
}

func (p *DBPostgre) LoadLinkPair(links T.DBlinksDTO) T.DBlinksDTO { // linkshort
	ctx := context.Background()
	query := "SELECT slink, llink FROM shortlink WHERE slink = $1"
	var tag1, tag2 string
	err := p.conn.QueryRow(ctx, query, links.Short).Scan(&tag1, &tag2)
	if err != nil {
		p.log.LogError(err, "(DBPostgre).LoadLinkPair(): SELECT error")
		return T.DBlinksDTO{}
	} else {
		p.log.LogDebug("(DBPostgre).LoadLinkPair(): %s, %s", tag1, tag2)
		return T.DBlinksDTO{Short: tag1, Long: tag2}
	}
}

func (p *DBPostgre) LoadAllLinkPairs() []T.DBlinksDTO {
	ctx := context.Background()
	query := "SELECT slink, llink FROM shortlink"
	rows, errRows := p.conn.Query(ctx, query)
	if errRows != nil {
		p.log.LogError(errRows, "(DBPostgre).LoadAllLinkPairs(): Query() SELECT error")
		return []T.DBlinksDTO{}
	}
	var tag1, tag2 string
	res := make([]T.DBlinksDTO, 0, 8)
	for rows.Next() {
		if errScan := rows.Scan(&tag1, &tag2); errScan != nil {
			p.log.LogError(errScan, "(DBPostgre).LoadAllLinkPairs(): Scan() SELECT error")
			return []T.DBlinksDTO{}
		}
		res = append(res, T.DBlinksDTO{Short: tag1, Long: tag2})
		p.log.LogDebug("(DBPostgre).LoadAllLinkPairs(): %s, %s", tag1, tag2)
	}
	return res
}

func (p *DBPostgre) Migration() {
	ctx := context.Background()

	query := "DROP TABLE IF EXISTS shortlink"
	tag, err := p.conn.Exec(ctx, query)
	if err != nil {
		p.log.LogError(err, "(DBPostgre).Migration(): DROP error")
	} else {
		p.log.LogDebug("(DBPostgre).Migration(): %s", tag)
	}

	query = "CREATE TABLE IF NOT EXISTS shortlink (slink TEXT PRIMARY KEY, llink TEXT NOT NULL, CHECK (llink <> ''))"
	tag, err = p.conn.Exec(ctx, query)
	if err != nil {
		p.log.LogError(err, "(DBPostgre).Migration(): CREATE error")
	} else {
		p.log.LogDebug("(DBPostgre).Migration(): %s", tag)
	}

	query = "INSERT INTO shortlink VALUES ('5clp60', 'http://lib.ru'); INSERT INTO shortlink VALUES ('dhiu79', 'http://google.ru');"
	tag, err = p.conn.Exec(ctx, query)
	if err != nil {
		p.log.LogError(err, "(DBPostgre).Migration(): INSERT error")
	} else {
		p.log.LogDebug("(DBPostgre).Migration(): %s", tag)
	}
}

func (p *DBPostgre) Connect() func(e error) {
	if p.cfg.SL_DB_PORT == "" {
		p.cfg.SL_DB_PORT = ":5432"
	}
	// pgURI := "postgres://login:pass@localhost:5432/database_name"
	pgURI := "postgres://" + p.cfg.SL_DB_LOGIN + ":" + p.cfg.SL_DB_PASS + "@" + p.cfg.SL_DB_IP + p.cfg.SL_DB_PORT + "/" + p.cfg.SL_DB_DBNAME
	pgpool, err := pgxpool.New(context.Background(), pgURI)
	if err != nil {
		p.log.LogError(err, "(DBPostgre).Connect(): unable to connect to postgres db: "+pgURI)
		return func(e error) {}
	} else {
		p.conn = pgpool
		p.log.LogInfo("postgres db connected: " + pgURI)
	}
	p.Migration()
	return func(e error) {
		pgpool.Close()
		if e != nil {
			p.log.LogError(e, "(DBPostgre).Connect(): postgres db disconnected with error")
		}
		p.log.LogInfo("postgres db disconnected")
	}
}
