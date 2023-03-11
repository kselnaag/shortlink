package adapters

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"shortlink/internal/ports"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
)

var _ ports.IHttpServer = (*HttpNetServer)(nil)

type HttpNetServer struct {
	servSL ports.IServShortLink
	hsrv   *gin.Engine
}

func NewHttpNetServer(servSL ports.IServShortLink) HttpNetServer {
	return HttpNetServer{
		servSL: servSL,
		hsrv:   gin.Default(),
	}
}

func (hns *HttpNetServer) Handle() ports.IHttpServer {
	type Message struct {
		IsResp bool
		Mode   string
		Body   string
	}
	hns.hsrv.GET("/check/ping", func(c *gin.Context) { // checks
		c.JSON(http.StatusOK, Message{true, "check", "pong"})
	})
	hns.hsrv.GET("/check/abs", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, Message{true, "check", "404 Not Found"})
	})
	hns.hsrv.GET("/check/panic", func(c *gin.Context) {
		panic(`{IsResp:true, Mode:check, Body:panic}`)
	})
	hns.hsrv.GET("/check/close", func(c *gin.Context) {
		hns.appClose()
		c.JSON(http.StatusOK, Message{true, "check", "close"})
	})
	//
	hns.hsrv.POST("/long", func(c *gin.Context) { // link short from link long
		body, _ := io.ReadAll(c.Request.Body)
		req := Message{}
		if err := json.Unmarshal(body, &req); (err != nil) || (req.IsResp) || (len(req.Body) == 0) {
			c.JSON(http.StatusBadRequest, Message{true, "app", req.Body})
			return
		}
		lp := hns.servSL.GetLinkShortFromLinkLong(req.Body)
		if !lp.IsValid() {
			c.JSON(http.StatusNotFound, Message{true, "app", req.Body})
			return
		}
		c.JSON(http.StatusOK, Message{true, "app", lp.Short()})
	})
	hns.hsrv.POST("/short", func(c *gin.Context) { // link long from link short
		body, _ := io.ReadAll(c.Request.Body)
		req := Message{}
		if err := json.Unmarshal(body, &req); (err != nil) || (req.IsResp) || (len(req.Body) != 6) {
			c.JSON(http.StatusBadRequest, Message{true, "app", req.Body})
			return
		}
		lp := hns.servSL.GetLinkLongFromLinkShort(req.Body)
		if !lp.IsValid() {
			c.JSON(http.StatusNotFound, Message{true, "app", req.Body})
			return
		}
		c.JSON(http.StatusOK, Message{true, "app", lp.Long()})
	})
	hns.hsrv.POST("/", func(c *gin.Context) { // save link pair
		body, _ := io.ReadAll(c.Request.Body)
		req := Message{}
		if err := json.Unmarshal(body, &req); (err != nil) || (req.IsResp) || (len(req.Body) == 0) {
			c.JSON(http.StatusBadRequest, Message{true, "app", req.Body})
			return
		}
		lp := hns.servSL.SetLinkPairFromLinkLong(req.Body)
		if !lp.IsValid() {
			c.JSON(http.StatusInternalServerError, Message{true, "app", req.Body})
			return
		}
		c.JSON(http.StatusOK, Message{true, "app", lp.Short()})
	})
	hns.hsrv.GET("/allpairs", func(c *gin.Context) { // all pairs
		strarr := []string{}
		pairs := hns.servSL.GetAllLinkPairs()
		for _, el := range pairs {
			strarr = append(strarr, el.Short()+": "+el.Long())
		}
		c.JSON(http.StatusOK, Message{true, "app", strings.Join(strarr, "\n")})
	})
	hns.hsrv.GET("/:hash", func(c *gin.Context) { // redirect
		hash := c.Param("hash")
		if len(hash) != 6 {
			c.JSON(http.StatusBadRequest, Message{true, "app", hash})
			return
		}
		lp := hns.servSL.GetLinkLongFromLinkShort(hash)
		if !lp.IsValid() {
			c.JSON(http.StatusBadRequest, Message{true, "app", hash})
			return
		}
		c.Redirect(http.StatusMovedPermanently, lp.Long())
	})
	hns.hsrv.GET("/", func(c *gin.Context) { //  app send UI
		//
		//
		//
		//
		c.JSON(http.StatusOK, Message{true, "app", "httpUI"})
	})
	return hns
} // curl -i -X POST localhost:8080/short -H 'Content-Type: application/json' -H 'Accept: application/json' -d '{"IsResp":false,"Mode":"client","Body":"5clp60"}'

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
		proc.Signal(syscall.SIGINT)
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
