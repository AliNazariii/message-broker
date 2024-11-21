package grpc

import (
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"therealbroker/pkg/config"
)

func Serve(conf *config.Grpc, grpcServer *grpc.Server) {
	listener, err := net.Listen("tcp", conf.Address)
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	logrus.Info("Start listening on address: ", conf.Address)

	if err = grpcServer.Serve(listener); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
}
