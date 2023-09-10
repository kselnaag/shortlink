package main

import (
	"os"
	"os/signal"
	"runtime"
	a9n "shortlink/internal"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(1) // GOMAXPROCS=1
	// debug.SetGCPercent(100) // GOGC=100
	// debug.SetMemoryLimit(2 831 155 200) // GOMEMLIMIT=2700MiB

	app := a9n.NewA9n()
	appStop := app.Start()

	defer func() {
		if err := recover(); err != nil {
			appStop(err.(error))
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-sig

	appStop(nil)
}
