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
	log := adapters.NewLogZero(cfg.HTTP_IP+cfg.HTTP_PORT, cfg.APP_NAME)
	db := adapters.NewDBMock() // cfg.DB_IP, cfg.DB_PORT
	hcli := adapters.NewHttpNetClient()
	svcsl := services.NewSvcShortLink(&db, &hcli, &log)
	hsrv := adapters.NewHttpFastServer(&svcsl, &log)
	hsrvShutdown := hsrv.Run(cfg.HTTP_PORT)
	// interrupt exec
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-sig
	// gracefull shutdown
	hsrvShutdown()
}
