package main

import (
	"context"
	"os"
	"os/signal"
	"shortlink/internal/adapters"
	"shortlink/internal/services"
	"syscall"
	"time"
)

func main() {
	cfg := adapters.NewCfgEnv("config.env")
	log := adapters.NewLogZero(cfg.HTTP_IP+cfg.HTTP_PORT, "shortlink")
	db := adapters.NewDBMock() // cfg.DB_IP, cfg.DB_PORT
	hcli := adapters.NewHttpNetClient()
	svcsl := services.NewServShortLink(&db, &hcli, &log)
	hserv := adapters.NewHttpNetServer(&svcsl)
	server := hserv.Handle().Run(cfg.HTTP_PORT)
	log.LogInfo("shortlink server running")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-sig

	ctxSHD, cancelSHD := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelSHD()
	if err := server.Shutdown(ctxSHD); err != nil {
		log.LogError(err, "shortlink server gracefull shutdown error")
	}
	// shutdown all other systems here

	log.LogInfo("shortLink server closed")
}
