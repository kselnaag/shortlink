package adapterHTTP_test

import (
	"net/http"
	adapterCfg "shortlink/internal/adapters/cfg"
	adapterDB "shortlink/internal/adapters/db"
	adapterHTTP "shortlink/internal/adapters/http"
	adapterLog "shortlink/internal/adapters/log"
	"shortlink/internal/services"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHTTPServer(t *testing.T) {
	asrt := assert.New(t)
	defer func() {
		err := recover()
		asrt.Nil(err)
	}()

	t.Run("HTTPNet", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		cfg := adapterCfg.CfgEnv{
			SL_APP_NAME:  "testSL",
			SL_HTTP_IP:   "localhost",
			SL_HTTP_PORT: ":8081",
			SL_DB_IP:     "localhost",
			SL_DB_PORT:   ":1313",
		}
		log := adapterLog.NewLogZero(&cfg)
		db := adapterDB.NewDBMock(&cfg)
		hcli := adapterHTTP.NewHTTPClientMock()
		svcsl := services.NewSvcShortLink(&db, &hcli, &log)
		hsrv := adapterHTTP.NewHTTPServerNet(&svcsl, &log, &cfg)
		_ = hsrv.Run()
		he := httpexpect.WithConfig(httpexpect.Config{
			Client: &http.Client{
				Transport: httpexpect.NewBinder(hsrv.Engine()),
				Jar:       httpexpect.NewCookieJar(),
			},
			Reporter: httpexpect.NewAssertReporter(t),
		})
		he.GET("/").Expect().Status(http.StatusOK).ContentType("text/html", "utf-8")
		he.GET("/check/ping").Expect().Status(http.StatusOK).ContentType("application/json", "utf-8").
			JSON().Object().HasValue("IsResp", true).HasValue("Mode", "check").HasValue("Body", "pong")
		he.GET("/check/abs").Expect().Status(http.StatusNotFound).ContentType("application/json", "utf-8").
			JSON().Object().HasValue("IsResp", true).HasValue("Mode", "check").HasValue("Body", "404 Not Found")
		he.GET("/check/panic").Expect().Status(http.StatusInternalServerError)
	})

	t.Run("HTTPFast", func(t *testing.T) {
		cfg := adapterCfg.CfgEnv{
			SL_APP_NAME:  "testSL",
			SL_HTTP_IP:   "localhost",
			SL_HTTP_PORT: ":8082",
			SL_DB_IP:     "localhost",
			SL_DB_PORT:   ":1313",
		}
		log := adapterLog.NewLogZero(&cfg)
		db := adapterDB.NewDBMock(&cfg)
		hcli := adapterHTTP.NewHTTPClientMock()
		svcsl := services.NewSvcShortLink(&db, &hcli, &log)
		hsrv := adapterHTTP.NewHTTPServerFast(&svcsl, &log, &cfg)
		_ = hsrv.Run()
		he := httpexpect.WithConfig(httpexpect.Config{
			Client: &http.Client{
				Transport: httpexpect.NewFastBinder(hsrv.Engine().Handler()),
				Jar:       httpexpect.NewCookieJar(),
			},
			Reporter: httpexpect.NewAssertReporter(t),
		})
		he.GET("/").Expect().Status(http.StatusOK).ContentType("text/html", "")
		he.GET("/check/ping").Expect().Status(http.StatusOK).ContentType("application/json", "").
			JSON().Object().HasValue("IsResp", true).HasValue("Mode", "check").HasValue("Body", "pong")
		he.GET("/check/abs").Expect().Status(http.StatusNotFound).ContentType("application/json", "").
			JSON().Object().HasValue("IsResp", true).HasValue("Mode", "check").HasValue("Body", "404 Not Found")
		he.GET("/check/panic").Expect().Status(http.StatusInternalServerError)
	})
}
