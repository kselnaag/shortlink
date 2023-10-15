package adapterDB_test

import (
	adapterDB "shortlink/internal/adapter/db"
	adapterLog "shortlink/internal/adapter/log"
	T "shortlink/internal/apptype"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresDB(t *testing.T) {
	asrt := assert.New(t)
	defer func() {
		err := recover()
		asrt.Nil(err)
	}()

	t.Run("dbPostgres", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping _dbPostgre_ tests in short mode")
		}
		cfg := T.CfgEnv{
			SL_APP_NAME:  "shortlink",
			SL_LOG_LEVEL: "trace",
			SL_HTTP_IP:   "localhost",
			SL_HTTP_PORT: ":8080",
			SL_DB_MODE:   "postgres",
			SL_DB_IP:     "localhost",
			SL_DB_PORT:   "",
			SL_DB_LOGIN:  "login",
			SL_DB_PASS:   "password",
			SL_DB_DBNAME: "shortlink",
		}
		log := adapterLog.NewLogFprintf(&cfg)
		pg := adapterDB.NewDBPostgres(&cfg, &log)
		dbShutdown := pg.Connect()

		links := T.DBlinksDTO{Short: "abcd", Long: "efjh"}
		asrt.True(pg.SaveLinkPair(links))
		links1 := pg.LoadLinkPair(T.DBlinksDTO{Short: "5clp60", Long: ""})
		asrt.Equal(T.DBlinksDTO{Short: "5clp60", Long: "http://lib.ru"}, links1)
		links2 := pg.LoadLinkPair(T.DBlinksDTO{Short: "abcd", Long: ""})
		asrt.Equal(T.DBlinksDTO{Short: "abcd", Long: "efjh"}, links2)
		links3 := pg.LoadAllLinkPairs()
		asrt.Equal([]T.DBlinksDTO{{Short: "5clp60", Long: "http://lib.ru"}, {Short: "dhiu79", Long: "http://google.ru"}, {Short: "abcd", Long: "efjh"}}, links3)

		pg.Migration()
		dbShutdown(nil)
	})
}

func TestMongoDB(t *testing.T) {
	asrt := assert.New(t)
	defer func() {
		err := recover()
		asrt.Nil(err)
	}()

	t.Run("dbMongo", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping _dbMongo_ tests in short mode")
		}
		cfg := T.CfgEnv{
			SL_APP_NAME:  "shortlink",
			SL_LOG_LEVEL: "trace",
			SL_HTTP_IP:   "localhost",
			SL_HTTP_PORT: ":8080",
			SL_DB_MODE:   "mongo",
			SL_DB_IP:     "localhost",
			SL_DB_PORT:   "",
			SL_DB_LOGIN:  "login",
			SL_DB_PASS:   "password",
			SL_DB_DBNAME: "shortlink",
		}
		log := adapterLog.NewLogFprintf(&cfg)
		mg := adapterDB.NewDBMongo(&cfg, &log)
		dbShutdown := mg.Connect()

		links := T.DBlinksDTO{Short: "abcd", Long: "efjh"}
		asrt.True(mg.SaveLinkPair(links))
		links1 := mg.LoadLinkPair(T.DBlinksDTO{Short: "5clp60", Long: ""})
		asrt.Equal(T.DBlinksDTO{Short: "5clp60", Long: "http://lib.ru"}, links1)
		links2 := mg.LoadLinkPair(T.DBlinksDTO{Short: "abcd", Long: ""})
		asrt.Equal(T.DBlinksDTO{Short: "abcd", Long: "efjh"}, links2)
		links3 := mg.LoadAllLinkPairs()
		asrt.Equal([]T.DBlinksDTO{{Short: "5clp60", Long: "http://lib.ru"}, {Short: "dhiu79", Long: "http://google.ru"}, {Short: "abcd", Long: "efjh"}}, links3)

		mg.Migration()
		dbShutdown(nil)
	})
}

func TestRedis(t *testing.T) {
	asrt := assert.New(t)
	defer func() {
		err := recover()
		asrt.Nil(err)
	}()

	t.Run("dbRedis", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping _dbRedis_ tests in short mode")
		}
		cfg := T.CfgEnv{
			SL_APP_NAME:  "shortlink",
			SL_LOG_LEVEL: "trace",
			SL_HTTP_IP:   "localhost",
			SL_HTTP_PORT: ":8080",
			SL_DB_MODE:   "redis",
			SL_DB_IP:     "localhost",
			SL_DB_PORT:   "",
			SL_DB_LOGIN:  "login",
			SL_DB_PASS:   "password",
			SL_DB_DBNAME: "shortlink",
		}
		log := adapterLog.NewLogFprintf(&cfg)
		rd := adapterDB.NewDBRedis(&cfg, &log)
		dbShutdown := rd.Connect()

		links := T.DBlinksDTO{Short: "abcd", Long: "efjh"}
		asrt.True(rd.SaveLinkPair(links))
		links1 := rd.LoadLinkPair(T.DBlinksDTO{Short: "5clp60", Long: ""})
		asrt.Equal(T.DBlinksDTO{Short: "5clp60", Long: "http://lib.ru"}, links1)
		links2 := rd.LoadLinkPair(T.DBlinksDTO{Short: "abcd", Long: ""})
		asrt.Equal(T.DBlinksDTO{Short: "abcd", Long: "efjh"}, links2)
		links3 := rd.LoadAllLinkPairs()
		asrt.Equal([]T.DBlinksDTO{{Short: "abcd", Long: "efjh"}, {Short: "dhiu79", Long: "http://google.ru"}, {Short: "5clp60", Long: "http://lib.ru"}}, links3)

		rd.Migration()
		dbShutdown(nil)
	})
}

func TestTarantool(t *testing.T) {
	asrt := assert.New(t)
	defer func() {
		err := recover()
		asrt.Nil(err)
	}()

	t.Run("dbTarantool", func(t *testing.T) { // # &"C:\Program Files\Go\bin\go.exe" test -v -tags go_tarantool_ssl_disable -vet=off -count=1 -run ^TestTarantool$ shortlink/internal/adapter/db
		if testing.Short() {
			t.Skip("skipping _dbTarantool_ tests in short mode")
		}
		cfg := T.CfgEnv{
			SL_APP_NAME:  "shortlink",
			SL_LOG_LEVEL: "trace",
			SL_HTTP_IP:   "localhost",
			SL_HTTP_PORT: ":8080",
			SL_DB_MODE:   "tarantool",
			SL_DB_IP:     "localhost",
			SL_DB_PORT:   "",
			SL_DB_LOGIN:  "login",
			SL_DB_PASS:   "password",
			SL_DB_DBNAME: "shortlink",
		}
		log := adapterLog.NewLogFprintf(&cfg)
		tt := adapterDB.NewDBTarantool(&cfg, &log)
		dbShutdown := tt.Connect()

		links := T.DBlinksDTO{Short: "abcd", Long: "efjh"}
		asrt.True(tt.SaveLinkPair(links))
		links1 := tt.LoadLinkPair(T.DBlinksDTO{Short: "5clp60", Long: ""})
		asrt.Equal(T.DBlinksDTO{Short: "5clp60", Long: "http://lib.ru"}, links1)
		links2 := tt.LoadLinkPair(T.DBlinksDTO{Short: "abcd", Long: ""})
		asrt.Equal(T.DBlinksDTO{Short: "abcd", Long: "efjh"}, links2)
		links3 := tt.LoadAllLinkPairs()
		asrt.Equal([]T.DBlinksDTO{{Short: "5clp60", Long: "http://lib.ru"}, {Short: "abcd", Long: "efjh"}, {Short: "dhiu79", Long: "http://google.ru"}}, links3)

		tt.Migration()
		dbShutdown(nil)
	})
}
