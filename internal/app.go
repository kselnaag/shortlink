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

func NewApp() *App {
	cfg := adapterCfg.NewCfgEnv("shortlink.env")
	log := createLogger(cfg)
	db := createDatabase(cfg, log)
	ctrlDB := control.NewCtrlDB(db)
	hcli := adapterHTTP.NewHTTPClientNet()
	svcsl := service.NewSvcShortLink(ctrlDB, hcli, log)
	ctrlHTTP := control.NewCtrlHTTP(svcsl)
	hsrv := createHTTPServer(ctrlHTTP, log, cfg)
	return &App{
		hsrv: hsrv,
		db:   db,
		log:  log,
		cfg:  cfg,
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

func createLogger(cfg *T.CfgEnv) T.ILog {
	switch cfg.SL_LOG_MODE {
	case "zerolog":
		return adapterLog.NewLogZero(cfg)
	case "logslog":
		return adapterLog.NewLogSlog(cfg)
	default:
		return adapterLog.NewLogFprintf(cfg)
	}
}

func createDatabase(cfg *T.CfgEnv, log T.ILog) T.IDB {
	switch cfg.SL_DB_MODE {
	case "postgres":
		return adapterDB.NewDBPostgres(cfg, log)
	case "mongodb":
		return adapterDB.NewDBMongo(cfg, log)
	case "redis":
		return adapterDB.NewDBRedis(cfg, log)
	case "tarantool":
		return adapterDB.NewDBTarantool(cfg, log)
	default:
		return adapterDB.NewDBMock(cfg, log)
	}
}

func createHTTPServer(ctrl T.ICtrlHTTP, log T.ILog, cfg *T.CfgEnv) T.IHTTPServer {
	switch cfg.SL_HTTP_MODE {
	case "fiber":
		return adapterHTTP.NewHTTPServerFast(ctrl, log, cfg)
	default:
		return adapterHTTP.NewHTTPServerNet(ctrl, log, cfg)
	}
}
