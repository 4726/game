package app

import (
	"log"
	"net"

	"github.com/4726/game/services/matchmaking/history/config"
	"github.com/4726/game/services/matchmaking/history/pb"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

type Service struct {
	cfg        config.Config
	hs         *historyServer
	grpcServer *grpc.Server
}

func NewService(cfg config.Config) *Service {
	return &Service{cfg, nil, nil}
}

func (s *Service) Run() error {
	lis, err := net.Listen("tcp", ":14000")
	if err != nil {
		log.Fatal(err)
	}

	s.hs, err = newHistoryServer(s.cfg)
	if err != nil {
		return err
	}

	s.grpcServer = grpc.NewServer()
	pb.RegisterHistoryServer(s.grpcServer, s.hs)
	return s.grpcServer.Serve(lis)
}

func (s *Service) Close() {
	if s.hs != nil {
		s.hs.Close()
	}
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}
