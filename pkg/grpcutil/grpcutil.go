package grpcutil

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

// DefaultServer returns a new grpc server with a logging middleware and a metrics middleware
func DefaultServer(logEntry *logrus.Entry) *grpc.Server {
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Internal, "panic: %v", p)
	}
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}

	var opts []grpc.ServerOption
	opts = append(opts, grpc_middleware.WithUnaryServerChain(
		grpc_logrus.UnaryServerInterceptor(logEntry),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
	))
	opts = append(opts, grpc_middleware.WithStreamServerChain(
		grpc_logrus.StreamServerInterceptor(logEntry),
		grpc_prometheus.StreamServerInterceptor,
		grpc_recovery.StreamServerInterceptor(recoveryOpts...),
	))
	grpcServer := grpc.NewServer(opts...)

	return grpcServer
}

func DefaultServerTLS(logEntry *logrus.Entry, tlsCertPath, tlsKeyPath string) (*grpc.Server, error) {
	creds, err := credentials.NewServerTLSFromFile(tlsCertPath, tlsKeyPath)
	if err != nil {
		return nil, err
	}

	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Internal, "panic: %v", p)
	}
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}

	var opts []grpc.ServerOption
	opts = append(opts, grpc.Creds(creds))
	opts = append(opts, grpc_middleware.WithUnaryServerChain(
		grpc_logrus.UnaryServerInterceptor(logEntry),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
	))
	opts = append(opts, grpc_middleware.WithStreamServerChain(
		grpc_logrus.StreamServerInterceptor(logEntry),
		grpc_prometheus.StreamServerInterceptor,
		grpc_recovery.StreamServerInterceptor(recoveryOpts...),
	))
	grpcServer := grpc.NewServer(opts...)

	return grpcServer, nil
}
