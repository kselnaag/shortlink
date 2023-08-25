package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"shortlink/internal/adapters"
	"shortlink/internal/services"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(1)

	cfg := adapters.NewCfgEnv("config.env")
	log := adapters.NewLogZero(&cfg)
	db := adapters.NewDBMock(&cfg)
	hcli := adapters.NewHTTPClientNet()
	svcsl := services.NewSvcShortLink(&db, &hcli, &log)
	hsrv := adapters.NewHTTPServerNet(&svcsl, &log, &cfg)
	hsrvShutdown := hsrv.Run()

	defer func() {
		if err := recover(); err != nil {
			log.LogError(fmt.Errorf("%s", err), "ERROR: shortlink service stoped")
		}
	}()

	log.LogInfo("shortlink service started")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-sig

	hsrvShutdown()
	// DB disconnect

	log.LogInfo("shortlink service stoped")
}
