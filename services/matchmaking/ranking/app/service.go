package app

import (
	"fmt"
	"net"

	grpcutil "github.com/4726/game/pkg/grpcutil"
	"github.com/4726/game/pkg/metrics"
	"github.com/4726/game/services/matchmaking/ranking/config"
	"github.com/4726/game/services/matchmaking/ranking/pb"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

type Service struct {
	cfg           config.Config
	rs            *rankingServer
	grpcServer    *grpc.Server
	metricsServer *metrics.HTTP
}

func NewService(cfg config.Config) *Service {
	return &Service{cfg, nil, nil, nil}
}

//Run runs the service and blocks until an error occurs
func (s *Service) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.cfg.Port))
	if err != nil {
		return err
	}

	s.rs, err = newRankingServer(s.cfg)
	if err != nil {
		return err
	}

	s.grpcServer = grpcutil.DefaultServer(logEntry)
	pb.RegisterRankingServer(s.grpcServer, s.rs)
	grpc_prometheus.Register(s.grpcServer)

	s.metricsServer = metrics.NewHTTP()
	go s.metricsServer.Run(s.cfg.Metrics.Port)

	return s.grpcServer.Serve(lis)
}

//Close gracefully stops the service
func (s *Service) Close() {
	if s.rs != nil {
		s.rs.Close()
	}
	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}
	s.metricsServer.Close()
}