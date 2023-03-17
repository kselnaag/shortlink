package adapters

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"regexp"
	"shortlink/internal/ports"
	"shortlink/web"
	"strings"
	"syscall"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var _ ports.IHttpServer = (*HttpNetServer)(nil)

type HttpNetServer struct {
	servSL ports.IServShortLink
	hsrv   *gin.Engine
	fs     embed.FS
}

func NewHttpNetServer(servSL ports.IServShortLink) HttpNetServer {
	return HttpNetServer{
		servSL: servSL,
		hsrv:   gin.Default(),
		fs:     web.StaticFS,
	}
}

func (hns *HttpNetServer) Handle() ports.IHttpServer {
	type Message struct {
		IsResp bool
		Mode   string
		Body   string
	}
	isHash := regexp.MustCompile(`^[a-z0-9]{6}$`).MatchString
	headers := func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache")
	}
	hns.hsrv.Use(static.Serve("/", NewEmbedFolder(hns.fs, "data")))

	hns.hsrv.GET("/check/ping", func(c *gin.Context) { // checks
		headers(c)
		c.JSON(http.StatusOK, Message{true, "check", "pong"})
	})
	hns.hsrv.GET("/check/abs", func(c *gin.Context) {
		headers(c)
		c.JSON(http.StatusNotFound, Message{true, "check", "404 Not Found"})
	})
	hns.hsrv.GET("/check/panic", func(c *gin.Context) {
		panic(`{IsResp:true, Mode:check, Body:panic}`)
	})
	hns.hsrv.GET("/check/close", func(c *gin.Context) {
		hns.appClose()
		headers(c)
		c.JSON(http.StatusOK, Message{true, "check", "close"})
	})
	hns.hsrv.GET("/allpairs", func(c *gin.Context) { // all pairs
		strarr := []string{}
		pairs := hns.servSL.GetAllLinkPairs()
		for _, el := range pairs {
			strarr = append(strarr, el.Short()+": "+el.Long())
		}
		headers(c)
		c.JSON(http.StatusOK, Message{true, "200", strings.Join(strarr, "; ")})
	})

	hns.hsrv.POST("/long", func(c *gin.Context) { // link short from link long
		body, _ := io.ReadAll(c.Request.Body)
		req := Message{}
		if err := json.Unmarshal(body, &req); (err != nil) || (req.IsResp) || (len(req.Body) == 0) {
			c.JSON(http.StatusBadRequest, Message{true, "400", req.Body})
			return
		}
		lp := hns.servSL.GetLinkShortFromLinkLong(req.Body)
		if !lp.IsValid() {
			c.JSON(http.StatusNotFound, Message{true, "404", req.Body})
			return
		}
		headers(c)
		c.JSON(http.StatusOK, Message{true, "200", lp.Short()})
	})
	hns.hsrv.POST("/short", func(c *gin.Context) { // link long from link short
		body, _ := io.ReadAll(c.Request.Body)
		req := Message{}
		if err := json.Unmarshal(body, &req); (err != nil) || (req.IsResp) || (!isHash(req.Body)) {
			c.JSON(http.StatusBadRequest, Message{true, "400", req.Body})
			return
		}
		lp := hns.servSL.GetLinkLongFromLinkShort(req.Body)
		if !lp.IsValid() {
			c.JSON(http.StatusNotFound, Message{true, "404", req.Body})
			return
		}
		headers(c)
		c.JSON(http.StatusOK, Message{true, "200", lp.Long()})
	})
	hns.hsrv.POST("/save", func(c *gin.Context) { // save link pair
		body, _ := io.ReadAll(c.Request.Body)
		req := Message{}
		if err := json.Unmarshal(body, &req); (err != nil) || (req.IsResp) || (len(req.Body) == 0) {
			c.JSON(http.StatusBadRequest, Message{true, "400", req.Body})
			return
		}
		lp := hns.servSL.SetLinkPairFromLinkLong(req.Body)
		if !lp.IsValid() {
			c.JSON(http.StatusInternalServerError, Message{true, "500", req.Body})
			return
		}
		headers(c)
		c.JSON(http.StatusOK, Message{true, "200", lp.Short()})
	})
	hns.hsrv.GET("/:hash", func(c *gin.Context) { // redirect
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
		headers(c)
		c.Redirect(http.StatusMovedPermanently, lp.Long())
	})
	/* hns.hsrv.GET("/favicon.png", func(c *gin.Context) {
		c.FileFromFS("data/favicon.png", http.FS(hns.fs))
	})*/
	hns.hsrv.GET("/", func(c *gin.Context) {
		data, err := hns.fs.ReadFile("data/index.html")
		if err != nil {
			c.JSON(http.StatusInternalServerError, Message{true, "500", "index.html not loaded"})
			return
		}
		headers(c)
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, string(data))
	})
	return hns
	// curl -i -X POST localhost:8080/save -H 'Content-Type: application/json' -H 'Accept: application/json' -d '{"IsResp":false,"Mode":"client","Body":"http://lib.ru/PROZA/"}'
	// Cache-Control: no-cache | Content-Type: text/html; charset=utf-8
	// (5clp60)http://lib.ru (dhiu79)http://google.ru (8b4s29)http://lib.ru/PROZA/
}

func (hns *HttpNetServer) Run(port string) *http.Server {
	srv := &http.Server{
		Addr:    port,
		Handler: hns.hsrv,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("server error: %s\n", err.Error())
		}
	}()
	return srv
}

func (hns *HttpNetServer) appClose() {
	if proc, err := os.FindProcess(syscall.Getpid()); err != nil {
		panic("pid not found: " + err.Error())
	} else {
		if err := proc.Signal(syscall.SIGINT); err != nil {
			panic("signar not sent: " + err.Error())
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

func NewEmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	subFS, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(subFS),
	}
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	trimed := strings.TrimPrefix(path, prefix)
	_, err := e.Open(trimed)
	return err == nil
}
