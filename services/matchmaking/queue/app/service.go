package app

import (
	"log"
	"net"

	"github.com/4726/game/services/matchmaking/queue/config"
	"github.com/4726/game/services/matchmaking/queue/pb"
	_ "github.com/go-sql-driver/mysql"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"google.golang.org/grpc"
)

type Service struct {
	cfg        config.Config
	qs         *queueServer
	grpcServer *grpc.Server
}

func NewService(cfg config.Config) *Service {
	return &Service{cfg, nil, nil}
}

//Run runs the service and blocks until an error occurs
func (s *Service) Run() error {
	lis, err := net.Listen("tcp", ":14000")
	if err != nil {
		log.Fatal(err)
	}

	s.qs = newQueueServer(s.cfg)

	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(grpc_logrus.UnaryServerInterceptor(logEntry)))
	opts = append(opts, grpc.StreamInterceptor(grpc_logrus.StreamServerInterceptor(logEntry)))
	s.grpcServer = grpc.NewServer(opts...)
	pb.RegisterQueueServer(s.grpcServer, s.qs)
	return s.grpcServer.Serve(lis)
}

//Close gracefully stops the service
func (s *Service) Close() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}
