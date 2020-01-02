package main

import (
	"log"
	"net"
	"github.com/4726/game/services/matchmaking/pb"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":14000")
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, )
	go func() {
		sig := <-c
		s.GracefulStop()
	}()

	server := grpc.NewServer()
	service := NewQueueService(QueueServiceOptions{100, 10})
	pb.RegisterQueueServer(server, service)

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
	}
}