package adapterCfg

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"shortlink/internal/i7e"
	"strings"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

func NewCfgEnv(cfgname string) i7e.CfgEnv {
	cfg := i7e.CfgEnv{ // default env
		SL_APP_NAME:  "shortlink",
		SL_HTTP_IP:   "localhost",
		SL_HTTP_PORT: ":8080",
		SL_DB_IP:     "localhost",
		SL_DB_PORT:   ":3301",
	}
	if ip, err := ipFromInterfaces(); err != nil {
		logMessage("error", cfg.SL_HTTP_IP+":"+cfg.SL_HTTP_PORT, cfg.SL_APP_NAME, err.Error(), "Can not get IP interface")
	} else {
		cfg.SL_HTTP_IP = ip
		cfg.SL_DB_IP = ip
	}

	exec, err := os.Executable() // LoadExecutablePath
	if err != nil {
		logMessage("error", cfg.SL_HTTP_IP+":"+cfg.SL_HTTP_PORT, cfg.SL_APP_NAME, err.Error(), "Executable not found")
	}
	filename := filepath.Join(filepath.Dir(exec), cfgname)
	if err := godotenv.Load(filename); err == nil { // LoadConfFromFileToEnv
		mess := fmt.Sprintf("Load config from file: ./%s", cfgname)
		logMessage("info", cfg.SL_HTTP_IP+":"+cfg.SL_HTTP_PORT, cfg.SL_APP_NAME, "", mess)
	}
	if err := env.Parse(&cfg); err != nil { // LoadConfFromEnvToStruct
		logMessage("error", cfg.SL_HTTP_IP+":"+cfg.SL_HTTP_PORT, cfg.SL_APP_NAME, err.Error(), "Use default config, environment parsing failed")
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

func logMessage(lvl, host, svc, err, mess string) {
	timenow := time.Now().Format(time.RFC3339Nano)
	fmt.Fprintf(os.Stderr, "{\"L\":\"%s\",\"T\":\"%s\",\"H\":\"%s\",\"S\":\"%s\",\"E\":\"%s\",\"M\":\"%s\"}\n", lvl, timenow, host, svc, err, mess)
}
