package adapters

import (
	"context"
	"fmt"
	"net/http"
	"shortlink/internal/ports"
	"shortlink/internal/services"

	"github.com/gin-gonic/gin"
)

var _ ports.IHttpServer = (*HttpNetServer)(nil)

type HttpNetServer struct {
	ctx    *context.Context
	servSL services.IServShortLink
	hsrv   *gin.Engine
}

func NewHttpNetServer(ctx *context.Context, servSL services.IServShortLink) HttpNetServer {
	return HttpNetServer{
		ctx:    ctx,
		servSL: servSL,
		hsrv:   gin.Default(),
	}
}

func (hns *HttpNetServer) Handle() ports.IHttpServer {
	hns.hsrv.GET("/", func(c *gin.Context) { // send UI
		c.String(http.StatusOK, "Hello UI\n")
	})
	hns.hsrv.GET("/ping", func(c *gin.Context) { // healthcheck
		c.String(http.StatusOK, "pong\n")
	})
	hns.hsrv.GET("/:hash", func(c *gin.Context) { // redirect
		hash := c.Param("hash")
		//c.String(http.StatusOK, "HASH: %s\n", hash)
		// c.Redirect(http.StatusMovedPermanently, "http://www.google.com/")
		c.String(http.StatusOK, "HASH: %s\n", hash)
	})
	hns.hsrv.GET("/allpairs", func(c *gin.Context) { // send all pairs
		c.String(http.StatusOK, "allpairs\n")
	})
	hns.hsrv.POST("/linkshort", func(c *gin.Context) { // link short from link long
		c.String(http.StatusOK, "link short\n")
	})
	hns.hsrv.POST("/linklong", func(c *gin.Context) { // link long from link short
		c.String(http.StatusOK, "link long\n")
	})
	return hns
}

func (hns *HttpNetServer) Run(port string) *http.Server {
	srv := &http.Server{
		Addr:    port,
		Handler: hns.hsrv,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
		}
	}()
	return srv
}

/*
get redirect from short link to long link
get html UI
get health check
get ALL link pairs presented in db
search the short link if you have a long link
search the long link if you have a short link
*/
