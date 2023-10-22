package adapterHTTP

import (
	"embed"
	"net/http"
	"os"
	"shortlink/web"
	"time"

	T "shortlink/internal/apptype"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	rec "github.com/gofiber/fiber/v2/middleware/recover"
)

var _ T.IHTTPServer = (*HTTPServerFast)(nil)

type HTTPServerFast struct {
	ctrl T.ICtrlHTTP
	hsrv *fiber.App
	fs   embed.FS
	log  T.ILog
	cfg  *T.CfgEnv
}

func NewHTTPServerFast(ctrl T.ICtrlHTTP, log T.ILog, cfg *T.CfgEnv) *HTTPServerFast {
	fiberconf := fiber.Config{
		Prefork:           false,
		CaseSensitive:     true,
		StrictRouting:     false,
		EnablePrintRoutes: true,
		UnescapePath:      true,
		ReadTimeout:       10 * time.Second,
		ServerHeader:      "fiber",
		AppName:           "shortlink",
	}
	return &HTTPServerFast{
		ctrl: ctrl,
		hsrv: fiber.New(fiberconf),
		fs:   web.StaticFS,
		log:  log,
		cfg:  cfg,
	}
}

func (hfs *HTTPServerFast) handlers() {
	headers := func(c *fiber.Ctx) {
		c.Set("Cache-Control", "no-cache")
	}
	hfs.hsrv.Use(logger.New())
	hfs.hsrv.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(hfs.fs), // static file server
		PathPrefix: "data",
		Browse:     false,
	}))
	hfs.hsrv.Use(rec.New())
	hfs.hsrv.Use(pprof.New())

	hfs.hsrv.Get("/check/ping", func(c *fiber.Ctx) error {
		headers(c)
		return c.Status(fiber.StatusOK).JSON(T.HTTPMessageDTO{IsResp: true, Mode: "check", Body: "pong"})
	})
	hfs.hsrv.Get("/check/abs", func(c *fiber.Ctx) error {
		headers(c)
		return c.Status(fiber.StatusNotFound).JSON(T.HTTPMessageDTO{IsResp: true, Mode: "check", Body: "404 Not Found"})
	})
	hfs.hsrv.Get("/check/panic", func(c *fiber.Ctx) error {
		headers(c)
		panic(`{IsResp:true,Mode:check,Body:panic}`)
		// return c.Status(fiber.StatusInternalServerError).JSON(HTTPMessageDTO{IsResp:true, Mode:"check", Body:"panic"})
	})
	hfs.hsrv.Get("/check/close", func(c *fiber.Ctx) error {
		headers(c)
		hfs.appClose()
		return c.Status(fiber.StatusOK).JSON(T.HTTPMessageDTO{IsResp: true, Mode: "check", Body: "server close"})
	})
	hfs.hsrv.Get("/check/allpairs", func(c *fiber.Ctx) error {
		headers(c)
		all, err := hfs.ctrl.AllPairs()
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(T.HTTPMessageDTO{IsResp: true, Mode: "404", Body: err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(T.HTTPMessageDTO{IsResp: true, Mode: "200", Body: all})
	})
	////
	hfs.hsrv.Post("/long", func(c *fiber.Ctx) error { // link short from link long
		headers(c)
		body := c.Body()
		short, err := hfs.ctrl.Long(body)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(T.HTTPMessageDTO{IsResp: true, Mode: "404", Body: err.Error()})
		}
		return c.Status(fiber.StatusPartialContent).JSON(T.HTTPMessageDTO{IsResp: true, Mode: "206", Body: short})
	})
	hfs.hsrv.Post("/short", func(c *fiber.Ctx) error { // link long from link short
		headers(c)
		body := c.Body()
		long, err := hfs.ctrl.Short(body)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(T.HTTPMessageDTO{IsResp: true, Mode: "404", Body: err.Error()})
		}
		return c.Status(fiber.StatusPartialContent).JSON(T.HTTPMessageDTO{IsResp: true, Mode: "206", Body: long})
	})
	hfs.hsrv.Post("/save", func(c *fiber.Ctx) error { // save link pair
		headers(c)
		body := c.Body()
		short, err := hfs.ctrl.Save(body)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(T.HTTPMessageDTO{IsResp: true, Mode: "404", Body: err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(T.HTTPMessageDTO{IsResp: true, Mode: "201", Body: short})
	})
	hfs.hsrv.Get("/r/:hash", func(c *fiber.Ctx) error { // redirect
		headers(c)
		hash := c.Params("hash")
		long, err := hfs.ctrl.Hash(hash)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(T.HTTPMessageDTO{IsResp: true, Mode: "404", Body: err.Error()})
		}
		c.Set("Content-Type", "text/html")
		return c.Redirect(long, fiber.StatusFound)
	})
}

func (hfs *HTTPServerFast) appClose() {
	hfs.log.LogInfo("fasthttp server closed by appClose() handle")
	os.Exit(0)
}

func (hfs *HTTPServerFast) Engine() *fiber.App {
	return hfs.hsrv
}

func (hfs *HTTPServerFast) Run() func(e error) {
	hfs.handlers()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				hfs.log.LogPanic(err.(error), "HTTPServerFast panic")
			}
		}()
		if err := hfs.hsrv.Listen(hfs.cfg.SL_HTTP_PORT); err != nil {
			hfs.log.LogError(err, "Run(): fasthttp server process error (closed)")
			hfs.appClose()
		} else {
			hfs.log.LogInfo("fasthttp server closed")
		}
	}()
	hfs.log.LogInfo("fasthttp server opened")
	return func(e error) {
		if err := hfs.hsrv.ShutdownWithTimeout(5 * time.Second); err != nil {
			hfs.log.LogError(err, "Run(): fasthttp server graceful_shutdown error")
		}
		if e != nil {
			hfs.log.LogError(e, "Run(): fasthttp server shutdown with error")
		}
	}
}

/* import "github.com/goccy/go-json"

func main() {
    app := fiber.New(fiber.Config{
        JSONEncoder: json.Marshal,
        JSONDecoder: json.Unmarshal,
    })
*/
