package adapters

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"runtime/debug"
	"shortlink/internal/ports"
	"shortlink/web"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	rec "github.com/gofiber/fiber/v2/middleware/recover"
)

var _ ports.IHTTPServer = (*HTTPServerFast)(nil)

type HTTPServerFast struct {
	servSL ports.ISvcShortLink
	hsrv   *fiber.App
	fs     embed.FS
	log    ports.ILog
	cfg    *CfgEnv
}

func NewHTTPServerFast(servSL ports.ISvcShortLink, log ports.ILog, cfg *CfgEnv) HTTPServerFast {
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
	return HTTPServerFast{
		servSL: servSL,
		hsrv:   fiber.New(fiberconf),
		fs:     web.StaticFS,
		log:    log,
		cfg:    cfg,
	}
}

func (hfs *HTTPServerFast) handlers() {
	type Message struct {
		IsResp bool
		Mode   string
		Body   string
	}
	isHash := regexp.MustCompile(`^[a-z0-9]{6}$`).MatchString
	headers := func(c *fiber.Ctx) {
		c.Set("Cache-Control", "no-cache")
	}
	logpanic := func() {
		if err := recover(); err != nil {
			hfs.log.LogPanic("%v\n%v", fmt.Sprintf("%v", err), string(debug.Stack()))
			panic(err)
		}
	}
	hfs.hsrv.Use(logger.New())
	hfs.hsrv.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(hfs.fs), // static file server
		PathPrefix: "data",
		Browse:     false,
	}))
	hfs.hsrv.Use(rec.New())

	hfs.hsrv.Get("/check/ping", func(c *fiber.Ctx) error {
		defer logpanic()
		headers(c)
		return c.Status(fiber.StatusOK).JSON(Message{true, "check", "pong"})
	})
	hfs.hsrv.Get("/check/abs", func(c *fiber.Ctx) error {
		defer logpanic()
		headers(c)
		return c.Status(fiber.StatusNotFound).JSON(Message{true, "check", "404 Not Found"})
	})
	hfs.hsrv.Get("/check/panic", func(c *fiber.Ctx) error {
		defer logpanic()
		headers(c)
		panic(`{IsResp:true,Mode:check,Body:panic}`)
		// return c.Status(fiber.StatusInternalServerError).JSON(Message{true, "check", "panic"})
	})
	hfs.hsrv.Get("/check/close", func(c *fiber.Ctx) error {
		defer logpanic()
		headers(c)
		hfs.appClose()
		return c.Status(fiber.StatusOK).JSON(Message{true, "check", "server close"})
	})
	hfs.hsrv.Get("/allpairs", func(c *fiber.Ctx) error {
		defer logpanic()
		headers(c)
		strarr := []string{}
		pairs := hfs.servSL.GetAllLinkPairs()
		for _, el := range pairs {
			strarr = append(strarr, el.Short()+": "+el.Long())
		}
		return c.Status(fiber.StatusOK).JSON(Message{true, "200", strings.Join(strarr, "; ")})
	})
	////
	hfs.hsrv.Post("/long", func(c *fiber.Ctx) error { // link short from link long
		defer logpanic()
		headers(c)
		body := c.Body()
		req := Message{}
		if err := json.Unmarshal(body, &req); (err != nil) || (req.IsResp) || (req.Body == "") {
			return c.Status(fiber.StatusBadRequest).JSON(Message{true, "400", req.Body})
		}
		lp := hfs.servSL.GetLinkShortFromLinkLong(req.Body)
		if !lp.IsValid() {
			return c.Status(fiber.StatusNotFound).JSON(Message{true, "404", req.Body})
		}
		return c.Status(fiber.StatusPartialContent).JSON(Message{true, "206", lp.Short()})
	})
	hfs.hsrv.Post("/short", func(c *fiber.Ctx) error { // link long from link short
		defer logpanic()
		headers(c)
		body := c.Body()
		req := Message{}
		if err := json.Unmarshal(body, &req); (err != nil) || (req.IsResp) || (!isHash(req.Body)) {
			return c.Status(fiber.StatusBadRequest).JSON(Message{true, "400", req.Body})
		}
		lp := hfs.servSL.GetLinkLongFromLinkShort(req.Body)
		if !lp.IsValid() {
			return c.Status(fiber.StatusNotFound).JSON(Message{true, "404", req.Body})
		}
		return c.Status(fiber.StatusPartialContent).JSON(Message{true, "206", lp.Long()})
	})
	hfs.hsrv.Post("/save", func(c *fiber.Ctx) error { // save link pair
		defer logpanic()
		headers(c)
		body := c.Body()
		req := Message{}
		if err := json.Unmarshal(body, &req); (err != nil) || (req.IsResp) || (req.Body == "") {
			return c.Status(fiber.StatusBadRequest).JSON(Message{true, "400", req.Body})
		}
		lp := hfs.servSL.SetLinkPairFromLinkLong(req.Body)
		if !lp.IsValid() {
			return c.Status(fiber.StatusBadRequest).JSON(Message{true, "400", req.Body})
		}
		return c.Status(fiber.StatusCreated).JSON(Message{true, "201", lp.Short()})
	})
	hfs.hsrv.Get("/r/:hash", func(c *fiber.Ctx) error { // redirect
		defer logpanic()
		headers(c)
		hash := c.Params("hash")
		if isHash(hash) {
			lp := hfs.servSL.GetLinkLongFromLinkShort(hash)
			if lp.IsValid() {
				return c.Redirect(lp.Long())
			}
		}
		return c.Status(fiber.StatusBadRequest).JSON(Message{true, "400", hash})
	})
}

func (hfs *HTTPServerFast) appClose() {
	if proc, err := os.FindProcess(syscall.Getpid()); err != nil {
		hfs.log.LogError(err, "appClose(): pid not found")
	} else {
		if err := proc.Signal(syscall.SIGINT); err != nil {
			hfs.log.LogError(err, "appClose(): signar not sent")
		}
	}
}

func (hfs *HTTPServerFast) Engine() *fiber.App {
	return hfs.hsrv
}

func (hfs *HTTPServerFast) Run() func() {
	hfs.handlers()
	go func() {
		if err := hfs.hsrv.Listen(hfs.cfg.HTTP_PORT); err != nil {
			hfs.log.LogError(err, "Run(): fasthttp server process error (closed)")
			hfs.appClose()
		} else {
			hfs.log.LogInfo("fasthttp server closed")
		}
	}()
	hfs.log.LogInfo("fasthttp server starting")
	return func() {
		if err := hfs.hsrv.ShutdownWithTimeout(5 * time.Second); err != nil {
			hfs.log.LogError(err, "Run(): fasthttp server graceful shutdown error")
		}
	}
}
