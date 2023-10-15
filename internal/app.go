package app

import (
	adapterCfg "shortlink/internal/adapter/cfg"
	adapterDB "shortlink/internal/adapter/db"
	adapterHTTP "shortlink/internal/adapter/http"
	adapterLog "shortlink/internal/adapter/log"
	T "shortlink/internal/apptype"
	"shortlink/internal/control"
	"shortlink/internal/service"
)

type App struct {
	hsrv T.IHTTPServer
	db   T.IDB
	log  T.ILog
	cfg  *T.CfgEnv
}

func NewApp() App {
	cfg := adapterCfg.NewCfgEnv("shortlink.env")
	log := adapterLog.NewLogFprintf(&cfg)
	db := adapterDB.NewDBMock(&cfg, &log)
	ctrlDB := control.NewCtrlDB(&db)
	hcli := adapterHTTP.NewHTTPClientNet()
	svcsl := service.NewSvcShortLink(&ctrlDB, &hcli, &log)
	ctrlHTTP := control.NewCtrlHTTP(&svcsl)
	hsrv := adapterHTTP.NewHTTPServerNet(&ctrlHTTP, &log, &cfg)
	return App{
		hsrv: &hsrv,
		db:   &db,
		log:  &log,
		cfg:  &cfg,
	}
}

func (a *App) Start() func(err error) {
	dbShutdown := a.db.Connect()
	hsrvShutdown := a.hsrv.Run()
	a.log.LogInfo(a.cfg.SL_APP_NAME + " app started")
	return func(err error) {
		hsrvShutdown(err)
		dbShutdown(err)
		if err != nil {
			a.log.LogError(err, a.cfg.SL_APP_NAME+" app stoped with error")
		} else {
			a.log.LogInfo(a.cfg.SL_APP_NAME + " app stoped")
		}
	}
}
