package main

import (
	"fmt"

	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTP_PORT string `env:"HTTP_PORT" envDefault:"8080"`
	DB_IP     string `env:"DB_IP" envDefault:"127.0.0.1"`
	DB_PORT   string `env:"DB_PORT" envDefault:"1313"`
}

var config Config

func init() {
	filename := "./config/shortlink.env"
	if err := godotenv.Load(filename); err != nil { // LoadConfFromFileToEnv
		fmt.Printf("No config file found: %s.\nLoad config from ENV variables directly.\n%s\n", filename, err.Error())
	} else {
		fmt.Printf("Load config from file: %s\n", filename)
	}
	if err := env.Parse(config); err != nil { // LoadConfFromEnvToStruct
		fmt.Printf("Setting default config, config parsing failed:\n%s\n", err.Error())
	}
}

func main() {
	/*
		ctx, cancel := context.WithCancel(context.Background())
		var wg sync.WaitGroup
		//wg.Add(1)
		go func() {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
			CYCSIG:
			for {
				select {
				case <-ctx.Done():
					break CYCSIG
				case <-sig:
					cancel()
					break CYCSIG
				}
			}
		}()

		//cancel()
		wg.Wait()
	*/
}
