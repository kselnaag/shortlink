package app

import (
	"os"
	adapterCfg "shortlink/internal/adapter/cfg"
	adapterDB "shortlink/internal/adapter/db"
	adapterHTTP "shortlink/internal/adapter/http"
	adapterLog "shortlink/internal/adapter/log"
	"shortlink/internal/control"
	"shortlink/internal/service"
	"shortlink/internal/types"
)

type App struct {
	hsrv types.IHTTPServer
	db   types.Idb
	log  types.ILog
	cfg  *types.CfgEnv
}

func NewApp() App {
	cfg := adapterCfg.NewCfgEnv("shortlink.env")
	log := adapterLog.NewLogZero(&cfg)
	db := adapterDB.NewDBMock(&cfg)
	hcli := adapterHTTP.NewHTTPClientNet()
	svcsl := service.NewSvcShortLink(&db, &hcli, &log)
	ctrl := control.NewCtrlHTTP(&svcsl)
	hsrv := adapterHTTP.NewHTTPServerNet(&ctrl, &log, &cfg)
	return App{
		hsrv: &hsrv,
		db:   &db,
		log:  &log,
		cfg:  &cfg,
	}
}

func (a *App) Start() func(err error) {
	// dbShutdown := db.Connect(&log, &cfg)
	hsrvShutdown := a.hsrv.Run()
	a.log.LogInfo(a.cfg.SL_APP_NAME + " app started")
	return func(err error) {
		hsrvShutdown(err)
		// dbShutdown(err)
		if err != nil {
			a.log.LogError(err, a.cfg.SL_APP_NAME+" app stoped")
			os.Exit(1)
		}
		a.log.LogInfo(a.cfg.SL_APP_NAME + " app stoped")
	}
}
