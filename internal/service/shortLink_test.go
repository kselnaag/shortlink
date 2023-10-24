package service_test

import (
	"testing"

	adapterDB "shortlink/internal/adapter/db"
	adapterHTTP "shortlink/internal/adapter/http"
	adapterLog "shortlink/internal/adapter/log"
	T "shortlink/internal/apptype"
	"shortlink/internal/control"
	"shortlink/internal/model"
	"shortlink/internal/service"

	"github.com/stretchr/testify/assert"
)

func TestServices(t *testing.T) {
	asrt := assert.New(t)
	defer func() {
		err := recover()
		asrt.Nil(err)
	}()

	t.Run("ServShortLink", func(t *testing.T) {
		cfg := &T.CfgEnv{
			SL_APP_NAME:  "testSL",
			SL_HTTP_IP:   "localhost",
			SL_HTTP_PORT: ":8080",
			SL_DB_IP:     "localhost",
			SL_DB_PORT:   ":1313",
		}
		log := adapterLog.NewLogFprintf(cfg)
		db := adapterDB.NewDBMock(cfg, log)
		ctrlDB := control.NewCtrlDB(db)
		hcli := adapterHTTP.NewHTTPClientMock(log)
		nssl := service.NewSvcShortLink(ctrlDB, hcli, log)

		// models.LinkPair{Short: "5clp60", Long: "http://lib.ru"}, models.LinkPair{Short: "8as3rb", Long: "http://lib.ru/abs"}, ("dhiu79", "http://google.ru")
		asrt.Equal([]model.LinkPair{model.NewLinkPair("http://lib.ru"), model.NewLinkPair("http://google.ru")},
			nssl.GetAllLinkPairs())

		asrt.True(nssl.IsLinkLongHTTPValid("http://lib.ru/PROZA/"))
		asrt.False(nssl.IsLinkLongHTTPValid("http://lib.ru/abs"))

		asrt.True(nssl.SetLinkPairFromLinkLong("http://lib.ru/PROZA/").IsValid())
		asrt.False(nssl.SetLinkPairFromLinkLong("http://lib.ru/abs").IsValid())

		lp := nssl.GetLinkLongFromLinkShort("8b4s29")
		asrt.True(lp.IsValid())
		asrt.Equal("http://lib.ru/PROZA/", lp.Long())
		asrt.False(nssl.GetLinkLongFromLinkShort("8as3rb").IsValid())

		asrt.Equal([]model.LinkPair{model.NewLinkPair("http://lib.ru"), model.NewLinkPair("http://lib.ru/PROZA/"), model.NewLinkPair("http://google.ru")},
			nssl.GetAllLinkPairs())
	})

}
