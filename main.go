package main

import (
	"google.golang.org/grpc"

	"therealbroker/api"
	"therealbroker/api/proto/src/broker/api/proto"
	"therealbroker/internal/broker"
	"therealbroker/internal/config"
	"therealbroker/internal/repositories"
	grpcPkg "therealbroker/pkg/grpc"
	interceptor "therealbroker/pkg/grpc/interceptor"
	"therealbroker/pkg/logger"
	"therealbroker/pkg/postgresql"
	"therealbroker/pkg/prometheus"
)

func main() {
	conf := config.New("broker")

	logger.Configure(conf.Log.Level)

	prometheus.StartPrometheusServer(conf.Prometheus.Port)

	brokerPostgres := postgresql.NewDB(&conf.Postgres, &repositories.MessagePostgres{})
	messageRepo := repositories.NewMessagesPostgres(brokerPostgres)

	//brokerCassandra := cassandra.NewDB(&conf.Cassandra)
	//messageRepo := repositories.NewCasMessageRepo(brokerCassandra, &conf.Cassandra)

	module := broker.NewModule(messageRepo)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.RPCMetricsInterceptor),
		grpc.UnaryInterceptor(interceptor.ErrorLoggingInterceptor),
	)
	handler := api.New(module, conf)
	proto.RegisterBrokerServer(grpcServer, handler)
	grpcPkg.Serve(&conf.Grpc, grpcServer)
}
