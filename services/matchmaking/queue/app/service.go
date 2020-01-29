package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/4726/game/services/matchmaking/queue/config"
	"github.com/4726/game/services/matchmaking/queue/pb"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
)

type Service struct {
	cfg        config.Config
	qs         *queueServer
	grpcServer *grpc.Server
	metricsServer *http.Server
}

func NewService(cfg config.Config) *Service {
	return &Service{cfg, nil, nil, nil}
}

//Run runs the service and blocks until an error occurs
func (s *Service) Run() error {
	lis, err := net.Listen("tcp", ":14000")
	if err != nil {
		log.Fatal(err)
	}

	s.qs = newQueueServer(s.cfg)

	var opts []grpc.ServerOption

	opts = append(opts, grpc_middleware.WithUnaryServerChain(
		grpc_logrus.UnaryServerInterceptor(logEntry),
		grpc_prometheus.UnaryServerInterceptor,
	))
	opts = append(opts, grpc_middleware.WithStreamServerChain(
		grpc_logrus.StreamServerInterceptor(logEntry),
		grpc_prometheus.StreamServerInterceptor,
	))
	s.grpcServer = grpc.NewServer(opts...)
	pb.RegisterQueueServer(s.grpcServer, s.qs)
	grpc_prometheus.Register(s.grpcServer)

	go s.runMetricsServer(s.cfg.Metrics) //need to autorestart if crashes maybe?

	return s.grpcServer.Serve(lis)
}

//Close gracefully stops the service
func (s *Service) Close() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
	if s.metricsServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		s.metricsServer.Shutdown(ctx)
	}
}

func (s *Service) runMetricsServer(metricsCfg config.MetricsConfig) error {
	s.metricsServer = &http.Server{Addr: fmt.Sprintf(":%v", metricsCfg.Port)}
	http.Handle(metricsCfg.Route, promhttp.Handler())
	return s.metricsServer.ListenAndServe()
}