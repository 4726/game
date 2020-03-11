package app

import (
	"net"

	"github.com/4726/game/gameplay/5v5fps/pb"
	grpcutil "github.com/4726/game/pkg/grpcutil"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

type Service struct {
	cs         *gameServer
	grpcServer *grpc.Server
}

// NewService returns a new Service
func NewService() (*Service, error) {
	s := &Service{}

	s.cs = newGameServer([5]uint64{1, 2, 3, 4, 5}, [5]uint64{6, 7, 8, 9, 10})
	s.grpcServer = grpcutil.DefaultServer(logEntry)

	pb.RegisterGameServer(s.grpcServer, s.cs)
	grpc_prometheus.Register(s.grpcServer)

	return s, nil
}

//Run runs the service and blocks until an error occurs
func (s *Service) Run() error {
	lis, err := net.Listen("tcp", ":14000")
	if err != nil {
		return err
	}

	logEntry.Info("started service on port 14000")
	return s.grpcServer.Serve(lis)
}

//Close gracefully stops the service
func (s *Service) Close() {
	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}
}
