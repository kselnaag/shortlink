package main

import (
	"os"
	"os/signal"
	"runtime"
	adapterCfg "shortlink/internal/adapter/cfg"
	adapterDB "shortlink/internal/adapter/db"
	adapterHTTP "shortlink/internal/adapter/http"
	adapterLog "shortlink/internal/adapter/log"
	"shortlink/internal/service"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(1) // GOMAXPROCS=1
	// debug.SetGCPercent(100) // GOGC=100
	// debug.SetMemoryLimit(2 831 155 200) // GOMEMLIMIT=2700MiB

	cfg := adapterCfg.NewCfgEnv("config.env")
	log := adapterLog.NewLogZero(&cfg)
	db := adapterDB.NewDBMock(&cfg)
	// dbShutdown := db.Connect()
	hcli := adapterHTTP.NewHTTPClientNet()
	svcsl := service.NewSvcShortLink(&db, &hcli, &log)
	hsrv := adapterHTTP.NewHTTPServerNet(&svcsl, &log, &cfg)
	hsrvShutdown := hsrv.Run()

	defer func() {
		if err := recover(); err != nil {
			hsrvShutdown()
			// dbShutdown()
			log.LogError(err.(error), "PANIC: "+cfg.SL_APP_NAME+" app stoped")
			os.Exit(1)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	log.LogInfo(cfg.SL_APP_NAME + " app started")
	<-sig

	hsrvShutdown()
	// dbShutdown()
	log.LogInfo(cfg.SL_APP_NAME + " app stoped")
}
