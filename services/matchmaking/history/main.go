package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/4726/game/services/matchmaking/history/app"
	"github.com/4726/game/services/matchmaking/history/pb"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":14000")
	if err != nil {
		log.Fatal(err)
	}

	cfg := app.Config{
		app.DBConfig{"history", "collection"},
		app.NSQConfig{"127.0.0.1:4150", "matches", "db"},
	}

	server := grpc.NewServer()
	service, err := app.NewService(cfg)
	if err != nil {
		log.Fatal(err)
	}
	pb.RegisterHistoryServer(server, service)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	serveCh := make(chan error, 1)
	go func() {
		err := server.Serve(lis)
		serveCh <- err
	}()

	select {
	case err := <-serveCh:
		log.Fatal(err)
	case sig := <-c:
		log.Fatal(sig.String())
		server.GracefulStop()
	}
}
