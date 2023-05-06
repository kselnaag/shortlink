package adapters

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"regexp"
	"runtime/debug"
	"shortlink/internal/ports"
	"shortlink/web"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var _ ports.IHTTPServer = (*HTTPServerNet)(nil)

type HTTPServerNet struct {
	servSL ports.ISvcShortLink
	hsrv   *gin.Engine
	fs     embed.FS
	log    ports.ILog
	cfg    *CfgEnv
}

func NewHTTPServerNet(servSL ports.ISvcShortLink, log ports.ILog, cfg *CfgEnv) HTTPServerNet {
	return HTTPServerNet{
		servSL: servSL,
		hsrv:   gin.New(),
		fs:     web.StaticFS,
		log:    log,
		cfg:    cfg,
	}
}

func (hns *HTTPServerNet) handlers() {
	type Message struct {
		IsResp bool
		Mode   string
		Body   string
	}
	isHash := regexp.MustCompile(`^[a-z0-9]{6}$`).MatchString
	headers := func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache")
	}
	logpanic := func(c *gin.Context) {
		if err := recover(); err != nil {
			hns.log.LogPanic("%v\n%v", fmt.Sprintf("%v", err), string(debug.Stack()))
			c.JSON(http.StatusInternalServerError, Message{true, "check", "panic"})
		}
	}
	hns.hsrv.Use(gin.Logger())
	hns.hsrv.Use(static.Serve("/", NewEmbedFolder(hns.fs, "data", hns.log)))

	hns.hsrv.GET("/check/ping", func(c *gin.Context) { // checks
		defer logpanic(c)
		headers(c)
		c.JSON(http.StatusOK, Message{true, "check", "pong"})
	})
	hns.hsrv.GET("/check/abs", func(c *gin.Context) {
		defer logpanic(c)
		headers(c)
		c.JSON(http.StatusNotFound, Message{true, "check", "404 Not Found"})
	})
	hns.hsrv.GET("/check/panic", func(c *gin.Context) {
		defer logpanic(c)
		headers(c)
		panic(`{IsResp:true,Mode:check,Body:panic}`)
		// c.JSON(http.StatusInternalServerError, Message{true, "check", "panic"})
	})
	hns.hsrv.GET("/check/close", func(c *gin.Context) {
		defer logpanic(c)
		headers(c)
		hns.appClose()
		c.JSON(http.StatusOK, Message{true, "check", "server close"})
	})
	hns.hsrv.GET("/allpairs", func(c *gin.Context) {
		defer logpanic(c)
		headers(c)
		strarr := []string{}
		pairs := hns.servSL.GetAllLinkPairs()
		for _, el := range pairs {
			strarr = append(strarr, el.Short()+": "+el.Long())
		}
		c.JSON(http.StatusOK, Message{true, "200", strings.Join(strarr, "; ")})
	})

	hns.hsrv.POST("/long", func(c *gin.Context) { // link short from link long
		defer logpanic(c)
		headers(c)
		body, readerr := io.ReadAll(c.Request.Body)
		req := Message{}
		if err := json.Unmarshal(body, &req); (err != nil) || (readerr != nil) || (req.IsResp) || (req.Body == "") {
			c.JSON(http.StatusBadRequest, Message{true, "400", req.Body})
			return
		}
		lp := hns.servSL.GetLinkShortFromLinkLong(req.Body)
		if !lp.IsValid() {
			c.JSON(http.StatusNotFound, Message{true, "404", req.Body})
			return
		}
		c.JSON(http.StatusPartialContent, Message{true, "206", lp.Short()})
	})
	hns.hsrv.POST("/short", func(c *gin.Context) { // link long from link short
		defer logpanic(c)
		headers(c)
		body, readerr := io.ReadAll(c.Request.Body)
		req := Message{}
		if err := json.Unmarshal(body, &req); (err != nil) || (readerr != nil) || (req.IsResp) || (!isHash(req.Body)) {
			c.JSON(http.StatusBadRequest, Message{true, "400", req.Body})
			return
		}
		lp := hns.servSL.GetLinkLongFromLinkShort(req.Body)
		if !lp.IsValid() {
			c.JSON(http.StatusNotFound, Message{true, "404", req.Body})
			return
		}
		c.JSON(http.StatusPartialContent, Message{true, "206", lp.Long()})
	})
	hns.hsrv.POST("/save", func(c *gin.Context) { // save link pair
		defer logpanic(c)
		headers(c)
		body, readerr := io.ReadAll(c.Request.Body)
		req := Message{}
		if err := json.Unmarshal(body, &req); (err != nil) || (readerr != nil) || (req.IsResp) || (req.Body == "") {
			c.JSON(http.StatusBadRequest, Message{true, "400", req.Body})
			return
		}
		lp := hns.servSL.SetLinkPairFromLinkLong(req.Body)
		if !lp.IsValid() {
			c.JSON(http.StatusBadRequest, Message{true, "400", req.Body})
			return
		}
		c.JSON(http.StatusCreated, Message{true, "201", lp.Short()})
	})
	hns.hsrv.GET("/:hash", func(c *gin.Context) { // redirect
		defer logpanic(c)
		headers(c)
		hash := c.Param("hash")
		if !isHash(hash) {
			c.JSON(http.StatusBadRequest, Message{true, "400", hash})
			return
		}
		lp := hns.servSL.GetLinkLongFromLinkShort(hash)
		if !lp.IsValid() {
			c.JSON(http.StatusBadRequest, Message{true, "400", hash})
			return
		}
		c.Redirect(http.StatusMovedPermanently, lp.Long())
	})
	/* hns.hsrv.GET("/favicon.png", func(c *gin.Context) {
		c.FileFromFS("data/favicon.png", http.FS(hns.fs))
	}) */
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
	})
	// curl -i -X POST localhost:8080/save -H 'Content-Type: application/json' -H 'Accept: application/json' -d '{"IsResp":false,"Mode":"client","Body":"http://lib.ru/PROZA/"}'
	// Cache-Control: no-cache | Content-Type: text/html; charset=utf-8
	// (5clp60)http://lib.ru (dhiu79)http://google.ru (8b4s29)http://lib.ru/PROZA/
}

func (hns *HTTPServerNet) appClose() {
	if proc, err := os.FindProcess(syscall.Getpid()); err != nil {
		hns.log.LogError(err, "appClose(): pid not found")
	} else {
		if err := proc.Signal(syscall.SIGINT); err != nil {
			hns.log.LogError(err, "appClose(): signar not sent")
		}
	}
}

func (hns *HTTPServerNet) Engine() *gin.Engine {
	return hns.hsrv
}

func (hns *HTTPServerNet) Run() func() {
	hns.handlers()
	srv := &http.Server{
		Addr:              hns.cfg.HTTP_PORT,
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
	hns.log.LogInfo("net/http server starting")
	return func() {
		ctxSHD, cancelSHD := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelSHD()
		if err := srv.Shutdown(ctxSHD); err != nil {
			hns.log.LogError(err, "Run(): net/http server graceful shutdown error")
		}
	}
}

/*
get redirect from short link to long link
get html UI
get health check
get ALL link pairs presented in db
search the short link if you have a long link
search the long link if you have a short link
*/

type embedFileSystem struct {
	http.FileSystem
}

func NewEmbedFolder(fsEmbed embed.FS, targetPath string, log ports.ILog) static.ServeFileSystem {
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
	_, err := e.Open(trimed)
	return err == nil
}
