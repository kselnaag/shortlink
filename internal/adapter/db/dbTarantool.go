package adapterDB

import (
	T "shortlink/internal/apptype"

	"github.com/tarantool/go-tarantool/v2"
)

var _ T.IDB = (*DBTarantool)(nil)

type row struct {
	SLINK string
	LLINK string
}

type DBTarantool struct {
	cfg  *T.CfgEnv
	log  T.ILog
	conn *tarantool.Connection
}

func NewDBTarantool(cfg *T.CfgEnv, log T.ILog) *DBTarantool {
	return &DBTarantool{
		cfg: cfg,
		log: log,
	}
}

func (t *DBTarantool) SaveLinkPair(links T.DBlinksDTO) bool {
	query := "INSERT INTO shortlink VALUES (?, ?)"
	req := tarantool.NewExecuteRequest(query).Args([]interface{}{links.Short, links.Long})
	if resp, err := t.conn.Do(req).Get(); err != nil {
		t.log.LogError(err, "(DBTarantool).SaveLinkPair(): INSERT error ")
		return false
	} else {
		t.log.LogDebug("(DBTarantool).SaveLinkPair(): INSERT %v", resp.SQLInfo)
		return true
	}
}

func (t *DBTarantool) LoadLinkPair(links T.DBlinksDTO) T.DBlinksDTO { // linkshort
	var resp []row
	req := tarantool.NewExecuteRequest("SELECT slink, llink FROM shortlink WHERE slink=?").Args([]interface{}{links.Short})
	if err := t.conn.Do(req).GetTyped(&resp); err != nil {
		t.log.LogDebug("(DBTarantool).LoadLinkPair(): SELECT error: %s", err.Error())
		return T.DBlinksDTO{}
	} else {
		t.log.LogDebug("(DBTarantool).LoadLinkPair(): SELECT %v", resp)
		if len(resp) != 1 {
			t.log.LogDebug("(DBTarantool).LoadLinkPair(): Num of rows != 1")
			return T.DBlinksDTO{}
		}
		return T.DBlinksDTO{Short: resp[0].SLINK, Long: resp[0].LLINK}
	}
}

func (t *DBTarantool) LoadAllLinkPairs() []T.DBlinksDTO {
	resp := make([]row, 0, 8)
	arrRes := make([]T.DBlinksDTO, 0, 8)
	query := "SELECT slink, llink FROM shortlink"
	req := tarantool.NewExecuteRequest(query)
	if err := t.conn.Do(req).GetTyped(&resp); err != nil {
		t.log.LogError(err, "(DBTarantool).LoadAllLinkPairs(): SELECT error ")
		return []T.DBlinksDTO{}
	} else {
		for _, el := range resp {
			arrRes = append(arrRes, T.DBlinksDTO{Short: el.SLINK, Long: el.LLINK})
			t.log.LogDebug("(DBTarantool).LoadAllLinkPairs(): SELECT %s %s", el.SLINK, el.LLINK)
		}
		return arrRes
	}
}

func (t *DBTarantool) Migration() {
	query := "DROP TABLE IF EXISTS shortlink"
	req := tarantool.NewExecuteRequest(query)
	if resp, err := t.conn.Do(req).Get(); err != nil {
		t.log.LogError(err, "(DBTarantool).Migration(): DROP TABLE error ")
	} else {
		t.log.LogDebug("(DBTarantool).Migration(): DROP TABLE %v", resp.SQLInfo)
	}

	query = "CREATE TABLE IF NOT EXISTS shortlink (slink TEXT PRIMARY KEY, llink TEXT NOT NULL, CHECK (llink <> ''))"
	req = tarantool.NewExecuteRequest(query)
	if resp, err := t.conn.Do(req).Get(); err != nil {
		t.log.LogError(err, "(DBTarantool).Migration(): CREATE TABLE error ")
	} else {
		t.log.LogDebug("(DBTarantool).Migration(): CREATE TABLE %v", resp.SQLInfo)
	}

	query = "INSERT INTO shortlink VALUES ('5clp60', 'http://lib.ru')"
	req = tarantool.NewExecuteRequest(query)
	if resp, err := t.conn.Do(req).Get(); err != nil {
		t.log.LogError(err, "(DBTarantool).Migration(): tarantool db <<box.execute>> INSERT error")
	} else {
		t.log.LogDebug("(DBTarantool).Migration(): INSERT %v", resp.SQLInfo)
	}

	query = "INSERT INTO shortlink VALUES ('dhiu79', 'http://google.ru')"
	req = tarantool.NewExecuteRequest(query)
	if resp, err := t.conn.Do(req).Get(); err != nil {
		t.log.LogError(err, "(DBTarantool).Migration(): INSERT error ")
	} else {
		t.log.LogDebug("(DBTarantool).Migration(): INSERT %v", resp.SQLInfo)
	}
}

func (t *DBTarantool) Connect() func(e error) {
	if t.cfg.SL_DB_PORT == "" {
		t.cfg.SL_DB_PORT = ":3301"
	}
	// ttURI := "localhost:3301"
	ttURI := t.cfg.SL_DB_IP + t.cfg.SL_DB_PORT
	conn, err := tarantool.Connect(ttURI, tarantool.Opts{User: t.cfg.SL_DB_LOGIN, Pass: t.cfg.SL_DB_PASS})
	if err != nil {
		t.log.LogError(err, "(DBTarantool).Connect(): unable to connect to tarantool db: "+ttURI)
		return func(e error) {}
	} else {
		t.conn = conn
		t.log.LogInfo("tarantool db connected: " + ttURI)
	}
	t.Migration()
	return func(e error) {
		if err := t.conn.CloseGraceful(); err != nil {
			t.log.LogError(err, "(DBTarantool).Connect(): tarantool db disconnection error")
		}
		if e != nil {
			t.log.LogError(e, "(DBTarantool).Connect(): tarantool db shutdown with error")
		}
		t.conn = nil
		t.log.LogInfo("tarantool db disconnected")
	}
}
