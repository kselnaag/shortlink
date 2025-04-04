package adapterHTTP_test

import (
	"net/http"
	adapterDB "shortlink/internal/adapter/db"
	adapterHTTP "shortlink/internal/adapter/http"
	adapterLog "shortlink/internal/adapter/log"
	T "shortlink/internal/apptype"
	"shortlink/internal/control"
	"shortlink/internal/service"
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
		cfg := &T.CfgEnv{
			SL_APP_NAME:  "testSL",
			SL_HTTP_IP:   "localhost",
			SL_HTTP_PORT: ":8081",
			SL_DB_IP:     "localhost",
			SL_DB_PORT:   ":1313",
		}
		log := adapterLog.NewLogFprintf(cfg)
		db := adapterDB.NewDBMock(cfg, log)
		ctrlDB := control.NewCtrlDB(db)
		hcli := adapterHTTP.NewHTTPClientMock(log)
		svcsl := service.NewSvcShortLink(ctrlDB, hcli, log)
		ctrlHTTP := control.NewCtrlHTTP(svcsl)
		hsrv := adapterHTTP.NewHTTPServerNet(ctrlHTTP, log, cfg)
		hsrvClose := hsrv.Run()
		he := httpexpect.WithConfig(httpexpect.Config{
			Client: &http.Client{
				Transport: httpexpect.NewBinder(hsrv.Engine()),
				Jar:       httpexpect.NewCookieJar(),
			},
			Reporter: httpexpect.NewAssertReporter(t),
		})

		he.GET("/").Expect().Status(http.StatusOK).HasContentType("text/html", "utf-8")
		he.GET("/check/ping").Expect().Status(http.StatusOK).HasContentType("application/json", "utf-8").
			JSON().Object().HasValue("IsResp", true).HasValue("Mode", "check").HasValue("Body", "pong")
		he.GET("/check/abs").Expect().Status(http.StatusNotFound).HasContentType("application/json", "utf-8").
			JSON().Object().HasValue("IsResp", true).HasValue("Mode", "check").HasValue("Body", "404 Not Found")
		asrt.Panics(func() { he.GET("/check/panic").Expect().Status(http.StatusInternalServerError) }, "HTTP Handle not panic")

		he.POST("/save").WithJSON(T.HTTPMessageDTO{IsResp: false, Mode: "save", Body: "http://lib.ru/PROZA/"}).Expect().
			Status(http.StatusCreated).HasContentType("application/json", "utf-8").
			JSON().Object().HasValue("IsResp", true).HasValue("Mode", "201").HasValue("Body", "8b4s29")
		he.POST("/short").WithJSON(T.HTTPMessageDTO{IsResp: false, Mode: "short", Body: "8b4s29"}).Expect().
			Status(http.StatusPartialContent).HasContentType("application/json", "utf-8").
			JSON().Object().HasValue("IsResp", true).HasValue("Mode", "206").HasValue("Body", "http://lib.ru/PROZA/")
		he.POST("/long").WithJSON(T.HTTPMessageDTO{IsResp: false, Mode: "long", Body: "http://lib.ru/PROZA/"}).Expect().
			Status(http.StatusPartialContent).HasContentType("application/json", "utf-8").
			JSON().Object().HasValue("IsResp", true).HasValue("Mode", "206").HasValue("Body", "8b4s29")
		he.GET("/check/allpairs").Expect().Status(http.StatusOK).HasContentType("application/json", "utf-8").
			JSON().Object().HasValue("IsResp", true).HasValue("Mode", "200").HasValue("Body", "5clp60: http://lib.ru; 8b4s29: http://lib.ru/PROZA/; dhiu79: http://google.ru")
		// he.GET("/r/8b4s29").Expect().Status(http.StatusFound).HasContentType("text/html", "utf-8")

		hsrvClose(nil)
	})

	t.Run("HTTPFast", func(t *testing.T) {
		cfg := &T.CfgEnv{
			SL_APP_NAME:  "testSL",
			SL_HTTP_IP:   "localhost",
			SL_HTTP_PORT: ":8082",
			SL_DB_IP:     "localhost",
			SL_DB_PORT:   ":1313",
		}
		log := adapterLog.NewLogFprintf(cfg)
		db := adapterDB.NewDBMock(cfg, log)
		ctrlDB := control.NewCtrlDB(db)
		hcli := adapterHTTP.NewHTTPClientMock(log)
		svcsl := service.NewSvcShortLink(ctrlDB, hcli, log)
		ctrlHTTP := control.NewCtrlHTTP(svcsl)
		hsrv := adapterHTTP.NewHTTPServerFast(ctrlHTTP, log, cfg)
		hsrvClose := hsrv.Run()
		he := httpexpect.WithConfig(httpexpect.Config{
			Client: &http.Client{
				Transport: httpexpect.NewFastBinder(hsrv.Engine().Handler()),
				Jar:       httpexpect.NewCookieJar(),
			},
			Reporter: httpexpect.NewAssertReporter(t),
		})
		he.GET("/").Expect().Status(http.StatusOK).HasContentType("text/html", "")
		he.GET("/check/ping").Expect().Status(http.StatusOK).HasContentType("application/json", "").
			JSON().Object().HasValue("IsResp", true).HasValue("Mode", "check").HasValue("Body", "pong")
		he.GET("/check/abs").Expect().Status(http.StatusNotFound).HasContentType("application/json", "").
			JSON().Object().HasValue("IsResp", true).HasValue("Mode", "check").HasValue("Body", "404 Not Found")
		he.GET("/check/panic").Expect().Status(http.StatusInternalServerError)

		he.POST("/save").WithJSON(T.HTTPMessageDTO{IsResp: false, Mode: "save", Body: "http://lib.ru/PROZA/"}).Expect().
			Status(http.StatusCreated).HasContentType("application/json", "").
			JSON().Object().HasValue("IsResp", true).HasValue("Mode", "201").HasValue("Body", "8b4s29")
		he.POST("/short").WithJSON(T.HTTPMessageDTO{IsResp: false, Mode: "short", Body: "8b4s29"}).Expect().
			Status(http.StatusPartialContent).HasContentType("application/json", "").
			JSON().Object().HasValue("IsResp", true).HasValue("Mode", "206").HasValue("Body", "http://lib.ru/PROZA/")
		he.POST("/long").WithJSON(T.HTTPMessageDTO{IsResp: false, Mode: "long", Body: "http://lib.ru/PROZA/"}).Expect().
			Status(http.StatusPartialContent).HasContentType("application/json", "").
			JSON().Object().HasValue("IsResp", true).HasValue("Mode", "206").HasValue("Body", "8b4s29")
		he.GET("/check/allpairs").Expect().Status(http.StatusOK).HasContentType("application/json", "").
			JSON().Object().HasValue("IsResp", true).HasValue("Mode", "200").HasValue("Body", "5clp60: http://lib.ru; 8b4s29: http://lib.ru/PROZA/; dhiu79: http://google.ru")
		// he.GET("/r/8b4s29").Expect().Status(http.StatusFound).ContentType("text/html", "utf-8")

		hsrvClose(nil)
	})
}
