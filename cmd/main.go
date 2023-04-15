package main

import (
	"os"
	"os/signal"
	"shortlink/internal/adapters"
	"shortlink/internal/services"
	"syscall"
)

func main() {
	//create and start all systems
	cfg := adapters.NewCfgEnv("config.env")
	log := adapters.NewLogZero(&cfg)
	db := adapters.NewDBMock(&cfg)
	hcli := adapters.NewHttpNetClient()
	svcsl := services.NewSvcShortLink(&db, &hcli, &log)
	hsrv := adapters.NewHttpNetServer(&svcsl, &log, &cfg)
	hsrvShutdown := hsrv.Run()
	// interrupt exec
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-sig
	// gracefull shutdown
	hsrvShutdown()
}
