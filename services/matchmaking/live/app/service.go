package app

import (
	"errors"
	"fmt"
	"net"

	grpcutil "github.com/4726/game/pkg/grpcutil"
	"github.com/4726/game/pkg/metrics"
	"github.com/4726/game/services/matchmaking/live/config"
	"github.com/4726/game/services/matchmaking/live/pb"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

type Service struct {
	cfg           config.Config
	ls            *liveServer
	grpcServer    *grpc.Server
	metricsServer *metrics.HTTP
}

// NewService returns a new Service
func NewService(cfg config.Config) (*Service, error) {
	s := &Service{}
	s.cfg = cfg

	var err error
	s.ls, err = newLiveServer(s.cfg)
	if err != nil {
		return nil, err
	}

	if s.cfg.TLS.CertPath != "" && s.cfg.TLS.KeyPath != "" {
		s.grpcServer, err = grpcutil.DefaultServerTLS(logEntry, s.cfg.TLS.CertPath, s.cfg.TLS.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("grpc error: %v", err)
		}
	} else {
		s.grpcServer = grpcutil.DefaultServer(logEntry)
	}

	pb.RegisterLiveServer(s.grpcServer, s.ls)
	grpc_prometheus.Register(s.grpcServer)

	s.metricsServer = metrics.NewHTTP()

	return s, nil
}

//Run runs the service and blocks until an error occurs
func (s *Service) Run() error {
	if s.metricsServer == nil || s.grpcServer == nil {
		return errors.New("service not setup, must call NewService()")
	}
	go s.metricsServer.Run(s.cfg.Metrics.Port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.cfg.Port))
	if err != nil {
		return err
	}

	logEntry.Info("started service on port: ", s.cfg.Port)
	return s.grpcServer.Serve(lis)
}

//Close gracefully stops the service
func (s *Service) Close() {
	if s.ls != nil {
		s.ls.Close()
	}
	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}
	if s.metricsServer != nil {
		s.metricsServer.Close()
	}
}
