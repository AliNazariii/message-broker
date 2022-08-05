package main

import (
	grpc "therealbroker/api"
	brokerModule "therealbroker/internal/broker"
	"therealbroker/internal/prometheus"
	"therealbroker/internal/repositories"
	"therealbroker/pkg/config"
	"therealbroker/pkg/database"
	"therealbroker/pkg/log"
)

// Main requirements:
// 1. All tests should be passed
// 2. Your logs should be accessible in Graylog
// 3. Basic prometheus metrics ( latency, throughput, etc. ) should be implemented
// 	  for every base functionality ( publish, subscribe etc. )

func main() {
	//conf := config.New("config.yaml", "broker")
	conf := config.New("", "broker")
	logger := log.NewLog(conf.Log.Level)
	metrics := prometheus.NewPrometheusServer(logger, conf.Prometheus.Port)

	brokerPostgres := database.NewPostgresDB(logger, &conf.Postgres, &repositories.PgMessage{})
	messageRepo := repositories.NewPgMessagesRepo(brokerPostgres, logger)

	//brokerCassandra := database.NewCassandraDB(logger, &conf.Cassandra)
	//messageRepo := repositories.NewCasMessageRepo(brokerCassandra, logger, &conf.Cassandra)

	broker := brokerModule.NewModule(logger, messageRepo)
	grpc.New(broker, logger, conf, metrics)
}
