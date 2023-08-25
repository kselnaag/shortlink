package adapters

import (
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type CfgEnv struct {
	SL_APP_NAME  string `env:"SL_APP_NAME"`
	SL_HTTP_IP   string `env:"SL_HTTP_IP"`
	SL_HTTP_PORT string `env:"SL_HTTP_PORT"`
	SL_DB_IP     string `env:"SL_DB_IP"`
	SL_DB_PORT   string `env:"SL_DB_PORT"`
}

func NewCfgEnv(cfgname string) CfgEnv {
	cfg := CfgEnv{ // default env
		SL_APP_NAME:  "shortlink",
		SL_HTTP_IP:   "localhost",
		SL_HTTP_PORT: ":8080",
		SL_DB_IP:     "localhost",
		SL_DB_PORT:   ":3301",
	}
	logCfg := NewLogZero(&cfg)
	if ip, err := IPFromInterfaces(); err != nil {
		logCfg.LogError(err, "Can not get IP interface")
	} else {
		cfg.SL_HTTP_IP = ip
		cfg.SL_DB_IP = ip
	}

	logCfg = NewLogZero(&cfg)
	exec, err := os.Executable() // LoadExecutablePath
	if err != nil {
		logCfg.LogError(err, "Executable not found")
	}
	filename := filepath.Join(filepath.Dir(exec), cfgname)
	if err := godotenv.Load(filename); err == nil { // LoadConfFromFileToEnv
		logCfg.LogInfo("Load config from file: ./%s", cfgname)
	}
	if err := env.Parse(&cfg); err != nil { // LoadConfFromEnvToStruct
		logCfg.LogError(err, "Use default config, environment parsing failed")
	}
	return cfg
}

func IPFromInterfaces() (string, error) {
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
