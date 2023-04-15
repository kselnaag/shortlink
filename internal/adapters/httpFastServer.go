package adapters

import (
	"embed"
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"shortlink/internal/ports"
	"shortlink/web"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var _ ports.IHttpServer = (*HttpFastServer)(nil)

type HttpFastServer struct {
	servSL ports.ISvcShortLink
	hsrv   *fiber.App
	fs     embed.FS
	log    ports.ILog
	cfg    *CfgEnv
}

func NewHttpFastServer(servSL ports.ISvcShortLink, log ports.ILog, cfg *CfgEnv) HttpFastServer {
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
	return HttpFastServer{
		servSL: servSL,
		hsrv:   fiber.New(fiberconf),
		fs:     web.StaticFS,
		log:    log,
		cfg:    cfg,
	}
}

func (hfs *HttpFastServer) handlers() {
	type Message struct {
		IsResp bool
		Mode   string
		Body   string
	}
	isHash := regexp.MustCompile(`^[a-z0-9]{6}$`).MatchString
	headers := func(c *fiber.Ctx) {
		c.Set("Cache-Control", "no-cache")
	}
	hfs.hsrv.Use(logger.New())
	hfs.hsrv.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(hfs.fs), // static file server
		PathPrefix: "data",
		Browse:     false,
	}))
	hfs.hsrv.Use(recover.New())

	hfs.hsrv.Get("/check/ping", func(c *fiber.Ctx) error {
		headers(c)
		return c.Status(fiber.StatusOK).JSON(Message{true, "check", "pong"})
	})
	hfs.hsrv.Get("/check/abs", func(c *fiber.Ctx) error {
		headers(c)
		return c.Status(fiber.StatusNotFound).JSON(Message{true, "check", "404 Not Found"})
	})
	hfs.hsrv.Get("/check/panic", func(c *fiber.Ctx) error {
		headers(c)
		panic(`{IsResp:true, Mode:check, Body:panic}`)
		//return c.Status(fiber.StatusInternalServerError).JSON(Message{true, "check", "panic"})
	})
	hfs.hsrv.Get("/check/close", func(c *fiber.Ctx) error {
		headers(c)
		hfs.appClose()
		return c.Status(fiber.StatusOK).JSON(Message{true, "check", "server close"})
	})
	hfs.hsrv.Get("/allpairs", func(c *fiber.Ctx) error {
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
		headers(c)
		body := c.Body()
		req := Message{}
		if err := json.Unmarshal(body, &req); (err != nil) || (req.IsResp) || (len(req.Body) == 0) {
			return c.Status(fiber.StatusBadRequest).JSON(Message{true, "400", req.Body})
		}
		lp := hfs.servSL.GetLinkShortFromLinkLong(req.Body)
		if !lp.IsValid() {
			return c.Status(fiber.StatusNotFound).JSON(Message{true, "404", req.Body})
		}
		return c.Status(fiber.StatusOK).JSON(Message{true, "200", lp.Short()})
	})
	hfs.hsrv.Post("/short", func(c *fiber.Ctx) error { // link long from link short
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
		return c.Status(fiber.StatusOK).JSON(Message{true, "200", lp.Long()})
	})
	hfs.hsrv.Post("/save", func(c *fiber.Ctx) error { // save link pair
		headers(c)
		body := c.Body()
		req := Message{}
		if err := json.Unmarshal(body, &req); (err != nil) || (req.IsResp) || (len(req.Body) == 0) {
			return c.Status(fiber.StatusBadRequest).JSON(Message{true, "400", req.Body})
		}
		lp := hfs.servSL.SetLinkPairFromLinkLong(req.Body)
		if !lp.IsValid() {
			return c.Status(fiber.StatusBadRequest).JSON(Message{true, "400", req.Body})
		}
		return c.Status(fiber.StatusOK).JSON(Message{true, "200", lp.Short()})
	})
	hfs.hsrv.Get("/:hash", func(c *fiber.Ctx) error { // redirect
		headers(c)
		hash := c.Params("hash")
		if !isHash(hash) {
			return c.Status(fiber.StatusBadRequest).JSON(Message{true, "400", hash})
		}
		lp := hfs.servSL.GetLinkLongFromLinkShort(hash)
		if !lp.IsValid() {
			return c.Status(fiber.StatusBadRequest).JSON(Message{true, "400", hash})
		}
		return c.Redirect(lp.Long())
	})
}

func (hfs *HttpFastServer) appClose() {
	if proc, err := os.FindProcess(syscall.Getpid()); err != nil {
		hfs.log.LogError(err, "appClose(): pid not found")
	} else {
		if err := proc.Signal(syscall.SIGINT); err != nil {
			hfs.log.LogError(err, "appClose(): signar not sent")
		}
	}
}

func (hfs *HttpFastServer) Run() func() {
	hfs.handlers()
	go func() {
		if err := hfs.hsrv.Listen(hfs.cfg.HTTP_PORT); err != nil {
			hfs.log.LogError(err, "Run(): fasthttp server process error")
		}
		hfs.log.LogInfo("fasthttp server closed")
	}()
	hfs.log.LogInfo("fasthttp server starting")
	return func() {
		if err := hfs.hsrv.ShutdownWithTimeout(5 * time.Second); err != nil {
			hfs.log.LogError(err, "Run(): fasthttp server gracefull shutdown error")
		}
	}
}
