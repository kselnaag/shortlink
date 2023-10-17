package adapterHTTP

import (
	"context"
	"embed"
	"io"
	"io/fs"
	"net/http"
	"os"
	T "shortlink/internal/apptype"
	"shortlink/web"
	"strings"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var _ T.IHTTPServer = (*HTTPServerNet)(nil)

type HTTPServerNet struct {
	ctrl T.ICtrlHTTP
	hsrv *gin.Engine
	fs   embed.FS
	log  T.ILog
	cfg  *T.CfgEnv
}

func NewHTTPServerNet(ctrl T.ICtrlHTTP, log T.ILog, cfg *T.CfgEnv) *HTTPServerNet {
	return &HTTPServerNet{
		ctrl: ctrl,
		hsrv: gin.New(),
		fs:   web.StaticFS,
		log:  log,
		cfg:  cfg,
	}
}

func (hns *HTTPServerNet) handlers() {
	headers := func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache")
	}
	hns.hsrv.Use(gin.Logger())
	hns.hsrv.Use(static.Serve("/", NewEmbedFolder(hns.fs, "data", hns.log)))
	/*
		[GIN-debug] GET    /debug/pprof/             --> github.com/gin-contrib/pprof.RouteRegister.WrapF.func1 (3 handlers)
		[GIN-debug] GET    /debug/pprof/cmdline      --> github.com/gin-contrib/pprof.RouteRegister.WrapF.func2 (3 handlers)
		[GIN-debug] GET    /debug/pprof/profile      --> github.com/gin-contrib/pprof.RouteRegister.WrapF.func3 (3 handlers)
		[GIN-debug] POST   /debug/pprof/symbol       --> github.com/gin-contrib/pprof.RouteRegister.WrapF.func4 (3 handlers)
		[GIN-debug] GET    /debug/pprof/symbol       --> github.com/gin-contrib/pprof.RouteRegister.WrapF.func5 (3 handlers)
		[GIN-debug] GET    /debug/pprof/trace        --> github.com/gin-contrib/pprof.RouteRegister.WrapF.func6 (3 handlers)
		[GIN-debug] GET    /debug/pprof/allocs       --> github.com/gin-contrib/pprof.RouteRegister.WrapH.func7 (3 handlers)
		[GIN-debug] GET    /debug/pprof/block        --> github.com/gin-contrib/pprof.RouteRegister.WrapH.func8 (3 handlers)
		[GIN-debug] GET    /debug/pprof/goroutine    --> github.com/gin-contrib/pprof.RouteRegister.WrapH.func9 (3 handlers)
		[GIN-debug] GET    /debug/pprof/heap         --> github.com/gin-contrib/pprof.RouteRegister.WrapH.func10 (3 handlers)
		[GIN-debug] GET    /debug/pprof/mutex        --> github.com/gin-contrib/pprof.RouteRegister.WrapH.func11 (3 handlers)
		[GIN-debug] GET    /debug/pprof/threadcreate --> github.com/gin-contrib/pprof.RouteRegister.WrapH.func12 (3 handlers)
	*/
	pprof.Register(hns.hsrv)

	hns.hsrv.GET("/check/ping", func(c *gin.Context) { // checks
		headers(c)
		c.JSON(http.StatusOK, T.HTTPMessageDTO{IsResp: true, Mode: "check", Body: "pong"})
	})
	hns.hsrv.GET("/check/abs", func(c *gin.Context) {
		headers(c)
		c.JSON(http.StatusNotFound, T.HTTPMessageDTO{IsResp: true, Mode: "check", Body: "404 Not Found"})
	})
	hns.hsrv.GET("/check/panic", func(c *gin.Context) {
		headers(c)
		panic(`{IsResp:true,Mode:check,Body:panic}`)
		// c.JSON(http.StatusInternalServerError, HTTPMessageDTO{true, "check", "panic"})
	})
	hns.hsrv.GET("/check/close", func(c *gin.Context) {
		headers(c)
		hns.appClose()
		c.JSON(http.StatusOK, T.HTTPMessageDTO{IsResp: true, Mode: "check", Body: "server close"})
	})
	hns.hsrv.GET("/check/allpairs", func(c *gin.Context) {
		headers(c)
		all, err := hns.ctrl.AllPairs()
		if err != nil {
			c.JSON(http.StatusNotFound, T.HTTPMessageDTO{IsResp: true, Mode: "404", Body: err.Error()})
			return
		}
		c.JSON(http.StatusOK, T.HTTPMessageDTO{IsResp: true, Mode: "200", Body: all})
	})

	hns.hsrv.POST("/long", func(c *gin.Context) { // link short from link long
		headers(c)
		body, readerr := io.ReadAll(c.Request.Body)
		if readerr != nil {
			c.JSON(http.StatusNotFound, T.HTTPMessageDTO{IsResp: true, Mode: "404", Body: readerr.Error()})
			return
		}
		short, err := hns.ctrl.Long(body)
		if err != nil {
			c.JSON(http.StatusNotFound, T.HTTPMessageDTO{IsResp: true, Mode: "404", Body: err.Error()})
			return
		}
		c.JSON(http.StatusPartialContent, T.HTTPMessageDTO{IsResp: true, Mode: "206", Body: short})
	})
	hns.hsrv.POST("/short", func(c *gin.Context) { // link long from link short
		headers(c)
		body, readerr := io.ReadAll(c.Request.Body)
		if readerr != nil {
			c.JSON(http.StatusNotFound, T.HTTPMessageDTO{IsResp: true, Mode: "404", Body: readerr.Error()})
			return
		}
		long, err := hns.ctrl.Short(body)
		if err != nil {
			c.JSON(http.StatusNotFound, T.HTTPMessageDTO{IsResp: true, Mode: "404", Body: err.Error()})
			return
		}
		c.JSON(http.StatusPartialContent, T.HTTPMessageDTO{IsResp: true, Mode: "206", Body: long})
	})
	hns.hsrv.POST("/save", func(c *gin.Context) { // save link pair
		headers(c)
		body, readerr := io.ReadAll(c.Request.Body)
		if readerr != nil {
			c.JSON(http.StatusNotFound, T.HTTPMessageDTO{IsResp: true, Mode: "404", Body: readerr.Error()})
			return
		}
		short, err := hns.ctrl.Save(body)
		if err != nil {
			c.JSON(http.StatusNotFound, T.HTTPMessageDTO{IsResp: true, Mode: "404", Body: err.Error()})
			return
		}
		c.JSON(http.StatusCreated, T.HTTPMessageDTO{IsResp: true, Mode: "201", Body: short})
	})
	hns.hsrv.GET("/r/:hash", func(c *gin.Context) { // redirect
		headers(c)
		hash := c.Param("hash")
		long, err := hns.ctrl.Hash(hash)
		if err != nil {
			c.JSON(http.StatusNotFound, T.HTTPMessageDTO{IsResp: true, Mode: "404", Body: err.Error()})
			return
		}
		c.Header("Content-Type", "text/html")
		c.Redirect(http.StatusFound, long)
	}) /*
		hns.hsrv.GET("/", func(c *gin.Context) {
			defer logpanic(c)
			headers(c)
			c.Header("Content-Type", "text/html; charset=utf-8")
			data, err := hns.fs.ReadFile("data/index.html")
			if err != nil {
				c.JSON(http.StatusInternalServerError, Message{true, "500", "embedFS: index.html not loaded"})
				return
			}
			c.String(http.StatusOK, string(data))
		}) */
	// curl -i -X POST localhost:8080/save -H 'Content-Type: application/json' -H 'Accept: application/json' -d '{"IsResp":false,"Mode":"client","Body":"http://lib.ru/PROZA/"}'
	// Cache-Control: no-cache | Content-Type: text/html; charset=utf-8
	// (5clp60)http://lib.ru (dhiu79)http://google.ru (8b4s29)http://lib.ru/PROZA/
}

func (hns *HTTPServerNet) appClose() {
	hns.log.LogInfo("net/http server closed by appClose() handle")
	os.Exit(0)
}

func (hns *HTTPServerNet) Engine() *gin.Engine {
	return hns.hsrv
}

func (hns *HTTPServerNet) Run() func(e error) {
	hns.handlers()
	srv := &http.Server{
		Addr:              hns.cfg.SL_HTTP_PORT,
		Handler:           hns.hsrv,
		ReadHeaderTimeout: 10 * time.Second,
	}
	go func() {
		err := srv.ListenAndServe()
		if (err != nil) && (err != http.ErrServerClosed) {
			hns.log.LogError(err, "Run(): net/http server process error (closed)")
			hns.appClose()
		}
		if err == http.ErrServerClosed {
			hns.log.LogInfo("net/http server closed")
		}
	}()
	hns.log.LogInfo("net/http server opened")
	return func(e error) {
		ctxSHD, cancelSHD := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelSHD()
		if err := srv.Shutdown(ctxSHD); err != nil {
			hns.log.LogError(err, "Run(): net/http server graceful_shutdown error")
		}
		if e != nil {
			hns.log.LogError(e, "Run(): net/http server shutdown with error")
		}
	}
}

type embedFileSystem struct {
	http.FileSystem
}

func NewEmbedFolder(fsEmbed embed.FS, targetPath string, log T.ILog) static.ServeFileSystem {
	subFS, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		log.LogError(err, "NewEmbedFolder(): embedFS error")
	}
	return embedFileSystem{
		FileSystem: http.FS(subFS),
	}
}

func (e embedFileSystem) Exists(prefix, path string) bool {
	trimed := strings.TrimPrefix(path, prefix)
	tlen := len(trimed)
	if tlen == 0 {
		trimed = "/"
		tlen++
	}
	if trimed[tlen-1] == '/' {
		trimed += "index.html"
	}
	_, err := e.Open(trimed)
	return err == nil
}
