package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/4726/game/services/matchmaking/history/app"
	"github.com/4726/game/services/matchmaking/history/config"
)

var configPath string

var usage = `
Usage: queue <command> [arguments]

Commands:
	-c, -config <file path> path to config file
	-h, -help	prints the usage string
`

func main() {
	flag.StringVar(&configPath, "config", "", "Path to config file")
	flag.StringVar(&configPath, "c", "", "Path to config file")
	flag.Usage = func() {
		fmt.Printf("%s\n", usage)
		os.Exit(0)
	}
	flag.Parse()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	service, err := app.NewService(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(c)

	serveCh := make(chan error, 1)
	go func() {
		err := service.Run()
		serveCh <- err
	}()

	select {
	case err := <-serveCh:
		log.Fatal(err)
	case sig := <-c:
		log.Fatal(sig.String())
	}
}
