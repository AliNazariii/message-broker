package grpc

import (
	"context"
	"google.golang.org/grpc/status"
	"therealbroker/pkg/metrics"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func RPCMetricsInterceptor(
	ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	method := info.FullMethod
	rpcStatus := "success"

	resp, err := handler(ctx, req)
	if err != nil {
		rpcStatus = "fail"
	}

	duration := time.Since(start).Seconds()

	metrics.MethodDuration.WithLabelValues(method, rpcStatus).Observe(duration)
	metrics.MethodCount.WithLabelValues(method, rpcStatus).Inc()

	return resp, err
}

func ErrorLoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	method := info.FullMethod

	resp, err := handler(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		logrus.WithError(err).Errorf("method %s failed with status: %s", method, st.Message())
	}

	return resp, err
}
