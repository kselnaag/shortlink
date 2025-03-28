package adapterCfg

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	adapterLog "shortlink/internal/adapter/log"
	T "shortlink/internal/apptype"
	"strings"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

var (
	ErrCfgEnv    = errors.New("CfgEnv error")
	ErrIPAddress = errors.New("IPaddress fail")
)

func NewCfgEnv(cfgname string) *T.CfgEnv {
	cfg := &T.CfgEnv{ // default env
		SL_APP_NAME:   "shortlink",
		SL_APP_PROTOS: ":http",
		SL_LOG_MODE:   "fprintf",
		SL_LOG_LEVEL:  "info",
		SL_HTTP_MODE:  "gin",
		SL_HTTP_IP:    "localhost",
		SL_HTTP_PORT:  ":8080",
		SL_DB_MODE:    "mock",
		SL_DB_IP:      "localhost",
		SL_DB_PORT:    "",
		SL_DB_LOGIN:   "login",
		SL_DB_PASS:    "password",
		SL_DB_DBNAME:  "shortlink",
	}
	log := adapterLog.NewLogFprintf(cfg)
	exec, err := os.Executable() // LoadExecutablePath
	if err != nil {
		log.LogError(err, "CfgEnv: os.Executable(): executable path not found")
	}
	filename := filepath.Join(filepath.Dir(exec), cfgname)
	log.LogDebug("CfgEnv exec path: %s", exec)
	NewCfgEnvFile(filename, cfg, log)
	return cfg
}

func NewCfgEnvFile(filename string, cfg *T.CfgEnv, log T.ILog) {
	if ip, err := ipFromInterfaces(); err != nil {
		log.LogError(err, "ipFromInterfaces(): can not get IP interface")
	} else {
		cfg.SL_HTTP_IP = ip
		log.LogDebug("(CfgEnv).ipFromInterfaces(): %s", ip)
	}
	if err := godotenv.Load(filename); err == nil { // LoadConfFromFileToEnv
		log.LogInfo("CfgEnv load config from file: %s", filename)
	}
	if err := env.Parse(cfg); err != nil { // LoadConfFromEnvToStruct
		log.LogWarn("env.Parse(): env vars parsing failed, use default config")
	}
}

func ipFromInterfaces() (string, error) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("%w: %w: %w", ErrCfgEnv, ErrIPAddress, err)
	}
	strarr := []string{}
	for _, addr := range addr {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				strarr = append(strarr, ipnet.IP.String())
			}
		}
	}
	return strings.Join(strarr, "; "), nil
}
