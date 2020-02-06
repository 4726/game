package grpcutil

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func DefaultServer(logEntry *logrus.Entry) *grpc.Server {
	var opts []grpc.ServerOption
	opts = append(opts, grpc_middleware.WithUnaryServerChain(
		grpc_logrus.UnaryServerInterceptor(logEntry),
		grpc_prometheus.UnaryServerInterceptor,
	))
	opts = append(opts, grpc_middleware.WithStreamServerChain(
		grpc_logrus.StreamServerInterceptor(logEntry),
		grpc_prometheus.StreamServerInterceptor,
	))
	grpcServer := grpc.NewServer(opts...)

	return grpcServer
}
