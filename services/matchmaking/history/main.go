package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"flag"

	"github.com/4726/game/services/matchmaking/history/app"
	"github.com/4726/game/services/matchmaking/history/pb"
	"google.golang.org/grpc"
	"github.com/4726/game/services/matchmaking/history/config"
)

var configPath string

func main() {
	flag.StringVar(&configPath, "config", "", "Path to config file")
	flag.StringVar(&configPath, "c", "", "Path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	service := app.NewService(cfg)
	defer service.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(c)

	serveCh := make(chan error, 1)
	go func() {
		err := server.Run()
		serveCh <- err
	}()

	select {
	case err := <-serveCh:
		log.Fatal(err)
	case sig := <-c:
		log.Fatal(sig.String())
	}
}
