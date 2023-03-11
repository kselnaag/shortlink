package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"shortlink/internal/adapters"
	"shortlink/internal/services"
	"strings"
	"syscall"
	"time"

	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTP_IP   string `env:"HTTP_IP"`
	HTTP_PORT string `env:"HTTP_PORT"`
	DB_IP     string `env:"DB_IP"`
	DB_PORT   string `env:"DB_PORT"`
}

func loadCfg() Config {
	var config = Config{ // default env
		HTTP_IP:   "localhost",
		HTTP_PORT: ":8080",
		DB_IP:     "localhost",
		DB_PORT:   ":1313",
	}
	exec, err := os.Executable() // LoadExecutablePath
	if err != nil {
		fmt.Printf("Executable not found: %s !\n", err.Error())
		os.Exit(1)
	}
	filename := filepath.Join(filepath.Dir(exec), "config.env")
	if err := godotenv.Load(filename); err == nil { // LoadConfFromFileToEnv
		fmt.Printf("Load config from file: %s .\n", filename)
	}
	if err := env.Parse(&config); err != nil { // LoadConfFromEnvToStruct
		fmt.Printf("Use default config, environment parsing failed !\n%s\n", err.Error())
	}
	return config
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

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// log := adapters.NewLogZero()
	// cfg := loadCfg(&log)
	db := adapters.NewMockDB()
	hcli := adapters.NewHttpNetClient()
	svcsl := services.NewServShortLink(&db, &hcli)
	hserv := adapters.NewHttpNetServer(&svcsl)
	server := hserv.Handle().Run(":8080")
	if ip, err := IpFromInterfaces(); err != nil {
		fmt.Printf("ShortLink server started at localhost%s\n%s\n", ":8080", err.Error())
	} else {
		fmt.Printf("ShortLink server started at %s%s\n", ip, ":8080")
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
SIGCYC:
	for {
		select {
		case <-ctx.Done():
			break SIGCYC
		case <-sig:
			break SIGCYC
		}
	}

	ctxSHD, cancelSHD := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelSHD()
	if err := server.Shutdown(ctxSHD); err != nil {
		fmt.Printf("SL server shutdown error: %s\n", err.Error())
	}
	// shutdown all other systems here

	fmt.Printf("\nShortLink server closed\n")
}
