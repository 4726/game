package grpcutil

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// DefaultServer returns a new grpc server with a logging middleware and a metrics middleware
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

func DefaultServerTLS(logEntry *logrus.Entry, tlsCertPath, tlsKeyPath string) (*grpc.Server, error) {
	creds, err := credentials.NewServerTLSFromFile(tlsCertPath, tlsKeyPath)
	if err != nil {
		return nil, err
	}

	var opts []grpc.ServerOption
	opts = append(opts, grpc.Creds(creds))
	opts = append(opts, grpc_middleware.WithUnaryServerChain(
		grpc_logrus.UnaryServerInterceptor(logEntry),
		grpc_prometheus.UnaryServerInterceptor,
	))
	opts = append(opts, grpc_middleware.WithStreamServerChain(
		grpc_logrus.StreamServerInterceptor(logEntry),
		grpc_prometheus.StreamServerInterceptor,
	))
	grpcServer := grpc.NewServer(opts...)

	return grpcServer, nil
}
