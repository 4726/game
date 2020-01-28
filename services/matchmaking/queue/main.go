package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/4726/game/services/matchmaking/queue/config"
	"github.com/4726/game/services/matchmaking/queue/pb"
	"google.golang.org/grpc"
	"github.com/4726/game/services/matchmaking/queue/app"
)

var configPath string

func main() {
	flag.StringVar(&configPath, "config", "", "Path to config file")
	flag.StringVar(&configPath, "c", "", "Path to config file")
	flag.Parse()

	lis, err := net.Listen("tcp", ":14000")
	if err != nil {
		log.Fatal(err)
	}

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
