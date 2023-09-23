package adapterCfg

import (
	"net"
	"os"
	"path/filepath"
	adapterLog "shortlink/internal/adapter/log"
	T "shortlink/internal/apptype"
	"strings"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

func NewCfgEnv(cfgname string) T.CfgEnv {
	cfg := T.CfgEnv{ // default env
		SL_APP_NAME:  "shortlink",
		SL_HTTP_IP:   "localhost",
		SL_HTTP_PORT: ":8080",
		SL_DB_IP:     "localhost",
		SL_DB_PORT:   ":3301",
		SL_LOG_LEVEL: "info",
	}
	log := adapterLog.NewLogFprintf(&cfg)
	if ip, err := ipFromInterfaces(); err != nil {
		log.LogError(err, "ipFromInterfaces(): can not get IP interface")
	} else {
		cfg.SL_HTTP_IP = ip
	}
	exec, err := os.Executable() // LoadExecutablePath
	if err != nil {
		log.LogError(err, "os.Executable(): executable path not found")
	}
	filename := filepath.Join(filepath.Dir(exec), cfgname)
	if err := godotenv.Load(filename); err == nil { // LoadConfFromFileToEnv
		log.LogInfo("Load config from file: ./%s", cfgname)
	}
	if err := env.Parse(&cfg); err != nil { // LoadConfFromEnvToStruct
		log.LogError(err, "env.Parse(): env vars parsing failed, use default config")
	}
	return cfg
}

func ipFromInterfaces() (string, error) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
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
