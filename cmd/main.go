package main

import (
	"os"
	"os/signal"
	"runtime"
	adapterCfg "shortlink/internal/adapters/cfg"
	adapterDB "shortlink/internal/adapters/db"
	adapterHTTP "shortlink/internal/adapters/http"
	adapterLog "shortlink/internal/adapters/log"
	"shortlink/internal/services"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(1)

	cfg := adapterCfg.NewCfgEnv("config.env")
	log := adapterLog.NewLogZero(&cfg)
	db := adapterDB.NewDBMock(&cfg)
	hcli := adapterHTTP.NewHTTPClientNet()
	svcsl := services.NewSvcShortLink(&db, &hcli, &log)
	hsrv := adapterHTTP.NewHTTPServerNet(&svcsl, &log, &cfg)
	hsrvShutdown := hsrv.Run()

	defer func() {
		if err := recover(); err != nil {
			log.LogError(err.(error), "PANIC: "+cfg.SL_APP_NAME+" app stoped")
			os.Exit(1)
		}
	}()

	log.LogInfo(cfg.SL_APP_NAME + " app started")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-sig

	hsrvShutdown()
	// DB disconnect

	log.LogInfo(cfg.SL_APP_NAME + " app stoped")
}
