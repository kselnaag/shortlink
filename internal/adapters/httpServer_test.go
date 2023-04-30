package adapters_test

import (
	"net/http"
	"shortlink/internal/adapters"
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
		cfg := adapters.CfgEnv{
			APP_NAME:  "testSL",
			HTTP_IP:   "localhost",
			HTTP_PORT: ":8081",
			DB_IP:     "localhost",
			DB_PORT:   ":1313",
		}
		log := adapters.NewLogZero(&cfg)
		db := adapters.NewDBMock(&cfg)
		hcli := adapters.NewHTTPClientMock()
		svcsl := services.NewSvcShortLink(&db, &hcli, &log)
		hsrv := adapters.NewHTTPServerNet(&svcsl, &log, &cfg)
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
		cfg := adapters.CfgEnv{
			APP_NAME:  "testSL",
			HTTP_IP:   "localhost",
			HTTP_PORT: ":8082",
			DB_IP:     "localhost",
			DB_PORT:   ":1313",
		}
		log := adapters.NewLogZero(&cfg)
		db := adapters.NewDBMock(&cfg)
		hcli := adapters.NewHTTPClientMock()
		svcsl := services.NewSvcShortLink(&db, &hcli, &log)
		hsrv := adapters.NewHTTPServerFast(&svcsl, &log, &cfg)
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
