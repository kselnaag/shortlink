package services_test

import (
	"testing"

	adaptersCfg "shortlink/internal/adapters/cfg"
	adapterDB "shortlink/internal/adapters/db"
	adapterHTTP "shortlink/internal/adapters/http"
	adapterLog "shortlink/internal/adapters/log"
	"shortlink/internal/models"
	"shortlink/internal/services"

	"github.com/stretchr/testify/assert"
)

func TestServices(t *testing.T) {
	asrt := assert.New(t)
	defer func() {
		err := recover()
		asrt.Nil(err)
	}()

	t.Run("ServShortLink", func(t *testing.T) {
		cfg := adaptersCfg.CfgEnv{
			SL_APP_NAME:  "testSL",
			SL_HTTP_IP:   "localhost",
			SL_HTTP_PORT: ":8080",
			SL_DB_IP:     "localhost",
			SL_DB_PORT:   ":1313",
		}
		log := adapterLog.NewLogZero(&cfg)
		db := adapterDB.NewDBMock(&cfg)
		hcli := adapterHTTP.NewHTTPClientMock()
		nssl := services.NewSvcShortLink(&db, &hcli, &log)

		// models.LinkPair{Short: "5clp60", Long: "http://lib.ru"}, models.LinkPair{Short: "8as3rb", Long: "http://lib.ru/abs"}, ("dhiu79", "http://google.ru")
		asrt.Equal([]models.LinkPair{models.NewLinkPair("http://lib.ru"), models.NewLinkPair("http://google.ru")},
			nssl.GetAllLinkPairs())

		asrt.True(nssl.IsLinkLongHTTPValid("http://lib.ru/PROZA/"))
		asrt.False(nssl.IsLinkLongHTTPValid("http://lib.ru/abs"))

		asrt.True(nssl.SetLinkPairFromLinkLong("http://lib.ru/PROZA/").IsValid())
		asrt.False(nssl.SetLinkPairFromLinkLong("http://lib.ru/abs").IsValid())

		lp := nssl.GetLinkLongFromLinkShort("8b4s29")
		asrt.True(lp.IsValid())
		asrt.Equal("http://lib.ru/PROZA/", lp.Long())
		asrt.False(nssl.GetLinkLongFromLinkShort("8as3rb").IsValid())

		asrt.Equal([]models.LinkPair{models.NewLinkPair("http://lib.ru"), models.NewLinkPair("http://lib.ru/PROZA/"), models.NewLinkPair("http://google.ru")},
			nssl.GetAllLinkPairs())
	})

}
