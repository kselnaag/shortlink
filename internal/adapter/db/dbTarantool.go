package adapterDB

import (
	"fmt"
	T "shortlink/internal/apptype"

	"github.com/tarantool/go-tarantool/v2"
)

var _ T.IDB = (*DBTarantool)(nil)

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
	query := fmt.Sprintf("INSERT INTO shortlink VALUES ('%s', '%s')", links.Short, links.Long)
	req := tarantool.NewCallRequest("box.execute").Args([]interface{}{query})
	if resp, err := t.conn.Do(req).Get(); err != nil {
		t.log.LogError(err, "(DBTarantool).SaveLinkPair(): tarantool db <<box.execute>> INSERT error")
		return false
	} else {
		t.log.LogDebug("(DBTarantool).SaveLinkPair(): INSERT %s", resp.Data...)
	}
	return true
}

func (t *DBTarantool) LoadLinkPair(links T.DBlinksDTO) T.DBlinksDTO { // linkshort
	query := fmt.Sprintf("SELECT slink, llink FROM shortlink WHERE slink = '%s'", links.Short)
	req := tarantool.NewCallRequest("box.execute").Args([]interface{}{query})
	resp, err := t.conn.Do(req).Get()
	if err != nil {
		t.log.LogDebug("(DBTarantool).LoadLinkPair(): tarantool db <<box.execute>> SELECT error: %s", err.Error())
		return T.DBlinksDTO{}
	} else {
		r := resp.Data[0].(map[interface{}]interface{})["rows"].([]interface{})[0]
		r1, ok1 := r.([]interface{})[0].(string)
		if !ok1 {
			t.log.LogDebug("(DBTarantool).LoadLinkPair(): fail r1 type assertion ")
		}
		r2, ok2 := r.([]interface{})[1].(string)
		if !ok2 {
			t.log.LogDebug("(DBTarantool).LoadLinkPair(): fail r2 type assertion")
		}
		t.log.LogDebug("(DBTarantool).LoadLinkPair(): SELECT %s %s", r1, r2)
		return T.DBlinksDTO{Short: r1, Long: r2}
	}
}

func (t *DBTarantool) LoadAllLinkPairs() []T.DBlinksDTO {
	arrRes := make([]T.DBlinksDTO, 0, 8)
	query := "SELECT slink, llink FROM shortlink"
	req := tarantool.NewCallRequest("box.execute").Args([]interface{}{query})
	resp, err := t.conn.Do(req).Get()
	if err != nil {
		t.log.LogError(err, "(DBTarantool).LoadAllLinkPairs(): tarantool db <<box.execute>> SELECT error")
		return []T.DBlinksDTO{}
	} else {
		r, ok := resp.Data[0].(map[interface{}]interface{})["rows"].([]interface{})
		if !ok {
			t.log.LogDebug("(DBTarantool).LoadAllLinkPairs(): fail r type assertion")
		}
		for _, el := range r {
			r1, ok1 := el.([]interface{})[0].(string)
			if !ok1 {
				t.log.LogDebug("(DBTarantool).LoadAllLinkPairs(): fail r1 type assertion")
			}
			r2, ok2 := el.([]interface{})[1].(string)
			if !ok2 {
				t.log.LogDebug("(DBTarantool).LoadAllLinkPairs(): fail r2 type assertion")
			}
			arrRes = append(arrRes, T.DBlinksDTO{Short: r1, Long: r2})
			t.log.LogDebug("(DBTarantool).LoadAllLinkPairs(): SELECT %s %s", r1, r2)
		}
		return arrRes
	}
}

func (t *DBTarantool) Migration() {
	req := tarantool.NewCallRequest("box.execute").Args([]interface{}{"DROP TABLE IF EXISTS shortlink"})
	if resp, err := t.conn.Do(req).Get(); err != nil {
		t.log.LogError(err, "(DBTarantool).Migration(): tarantool db <<box.execute>> DROP TABLE error")
	} else {
		t.log.LogDebug("(DBTarantool).Migration(): DROP TABLE %s", resp.Data...)
	}
	req = tarantool.NewCallRequest("box.execute").Args([]interface{}{"CREATE TABLE IF NOT EXISTS shortlink (slink TEXT PRIMARY KEY, llink TEXT NOT NULL, CHECK (llink <> ''))"})
	if resp, err := t.conn.Do(req).Get(); err != nil {
		t.log.LogError(err, "(DBTarantool).Migration(): tarantool db <<box.execute>> CREATE TABLE error")
	} else {
		t.log.LogDebug("(DBTarantool).Migration(): CREATE TABLE %s", resp.Data...)
	}

	req = tarantool.NewCallRequest("box.execute").Args([]interface{}{"INSERT INTO shortlink VALUES ('5clp60', 'http://lib.ru')"})
	if resp, err := t.conn.Do(req).Get(); err != nil {
		t.log.LogError(err, "(DBTarantool).Migration(): tarantool db <<box.execute>> INSERT error")
	} else {
		t.log.LogDebug("(DBTarantool).Migration(): INSERT %s", resp.Data...)
	}
	req = tarantool.NewCallRequest("box.execute").Args([]interface{}{"INSERT INTO shortlink VALUES ('dhiu79', 'http://google.ru')"})
	if resp, err := t.conn.Do(req).Get(); err != nil {
		t.log.LogError(err, "(DBTarantool).Migration(): tarantool db <<box.execute>> INSERT error")
	} else {
		t.log.LogDebug("(DBTarantool).Migration(): INSERT %s", resp.Data...)
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
