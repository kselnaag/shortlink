package adapterHTTP_test

import (
	"net/http"
	adapterDB "shortlink/internal/adapter/db"
	adapterHTTP "shortlink/internal/adapter/http"
	adapterLog "shortlink/internal/adapter/log"
	"shortlink/internal/control"
	"shortlink/internal/service"
	"shortlink/internal/types"
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
		cfg := types.CfgEnv{
			SL_APP_NAME:  "testSL",
			SL_HTTP_IP:   "localhost",
			SL_HTTP_PORT: ":8081",
			SL_DB_IP:     "localhost",
			SL_DB_PORT:   ":1313",
		}
		log := adapterLog.NewLogZero(&cfg)
		db := adapterDB.NewDBMock(&cfg)
		hcli := adapterHTTP.NewHTTPClientMock()
		svcsl := service.NewSvcShortLink(&db, &hcli, &log)
		ctrl := control.NewCtrlHTTP(&svcsl)
		hsrv := adapterHTTP.NewHTTPServerNet(&ctrl, &log, &cfg)
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
		asrt.Panics(func() { he.GET("/check/panic").Expect().Status(http.StatusInternalServerError) }, "HTTP Handle not panic")
	})

	t.Run("HTTPFast", func(t *testing.T) {
		cfg := types.CfgEnv{
			SL_APP_NAME:  "testSL",
			SL_HTTP_IP:   "localhost",
			SL_HTTP_PORT: ":8082",
			SL_DB_IP:     "localhost",
			SL_DB_PORT:   ":1313",
		}
		log := adapterLog.NewLogZero(&cfg)
		db := adapterDB.NewDBMock(&cfg)
		hcli := adapterHTTP.NewHTTPClientMock()
		svcsl := service.NewSvcShortLink(&db, &hcli, &log)
		ctrl := control.NewCtrlHTTP(&svcsl)
		hsrv := adapterHTTP.NewHTTPServerFast(&ctrl, &log, &cfg)
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
