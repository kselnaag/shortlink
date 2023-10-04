package adapterDB_test

import (
	adapterDB "shortlink/internal/adapter/db"
	adapterLog "shortlink/internal/adapter/log"
	T "shortlink/internal/apptype"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	asrt := assert.New(t)
	defer func() {
		err := recover()
		asrt.Nil(err)
	}()

	t.Run("dbPostgre", func(t *testing.T) {
		cfg := T.CfgEnv{
			SL_APP_NAME:  "shortlink",
			SL_LOG_LEVEL: "trace",
			SL_HTTP_IP:   "localhost",
			SL_HTTP_PORT: ":8080",
			SL_DB_IP:     "localhost",
			SL_DB_PORT:   ":5432",
			SL_DB_LOGIN:  "postgres",
			SL_DB_PASS:   "example",
			SL_DB_DBNAME: "postgres",
		}

		log := adapterLog.NewLogFprintf(&cfg)
		pg := adapterDB.NewDBPostgre(&cfg, &log)
		shutdown := pg.Connect()

		links := T.DBlinksDTO{Short: "abcd", Long: "efjh"}
		asrt.True(pg.SaveLinkPair(links))

		links1 := pg.LoadLinkPair(T.DBlinksDTO{Short: "5clp60", Long: ""})
		asrt.Equal(links1, T.DBlinksDTO{Short: "5clp60", Long: "http://lib.ru"})

		links2 := pg.LoadLinkPair(T.DBlinksDTO{Short: "abcd", Long: ""})
		asrt.Equal(links2, T.DBlinksDTO{Short: "abcd", Long: "efjh"})

		links3 := pg.LoadAllLinkPairs()
		asrt.Equal(links3, []T.DBlinksDTO{{Short: "5clp60", Long: "http://lib.ru"}, {Short: "dhiu79", Long: "http://google.ru"}, {Short: "abcd", Long: "efjh"}})

		pg.Migration()
		shutdown(nil)
	})
}
