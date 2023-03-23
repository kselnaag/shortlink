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
	APP_NAME  string `env:"APP_NAME"`
	HTTP_IP   string
	HTTP_PORT string `env:"HTTP_PORT"`
	DB_IP     string `env:"DB_IP"`
	DB_PORT   string `env:"DB_PORT"`
}

func NewCfgEnv(cfgname string) CfgEnv {
	cfg := CfgEnv{ // default env
		APP_NAME:  "shortlink",
		HTTP_IP:   "localhost",
		HTTP_PORT: ":8080",
		DB_IP:     "localhost",
		DB_PORT:   ":1313",
	}
	logCfg := LogZero{}
	if ip, err := IpFromInterfaces(); err != nil {
		logCfg = NewLogZero("localhost", "shortlink")
		logCfg.LogError(err, "Can not get IP interface")
	} else {
		cfg.HTTP_IP = ip
		logCfg = NewLogZero(cfg.HTTP_IP, "shortlink")
	}

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

func IpFromInterfaces() (string, error) {
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
